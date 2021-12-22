package connect

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	// ErrForbidden will be returned if the client doesn't have access to the
	// requested ressource.
	ErrForbidden = Error("forbidden")

	// ErrNotFound will be returned if the requested ressource could not be
	// found.
	ErrNotFound = Error("not found")

	// ErrBadRequest will be returned if Garmin returned a status code 400.
	ErrBadRequest = Error("bad request")

	// ErrNoCredentials will be returned if credentials are needed - but none
	// are set.
	ErrNoCredentials = Error("no credentials set")

	// ErrNotAuthenticated will be returned is the client is not
	// authenticated as required by the request. Remember to call
	// Authenticate().
	ErrNotAuthenticated = Error("client is not authenticated")

	// ErrWrongCredentials will be returned if the username and/or
	// password is not recognized by Garmin Connect.
	ErrWrongCredentials = Error("username and/or password not recognized")
)

const (
	// sessionCookieName is the magic session cookie name.
	sessionCookieName = "SESSIONID"

	// cflbCookieName is the cookie used by Cloudflare to pin the request
	// to a specific backend.
	cflbCookieName = "__cflb"
)

// Client can be used to access the unofficial Garmin Connect API.
type Client struct {
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	SessionID string         `json:"sessionID"`
	Profile   *SocialProfile `json:"socialProfile"`

	// LoadBalancerID is the load balancer ID set by Cloudflare in front of
	// Garmin Connect. This must be preserves across requests. A session key
	// is only valid with a corresponding loadbalancer key.
	LoadBalancerID string `json:"cflb"`

	client           *http.Client
	autoRenewSession bool
	debugLogger      Logger
	dumpWriter       io.Writer
}

// Option is the type to set options on the client.
type Option func(*Client)

// SessionID will set a predefined session ID. This can be useful for clients
// keeping state. A few HTTP roundtrips can be saved, if the session ID is
// reused. And some load would be taken of Garmin servers. This must be
// accompanied by LoadBalancerID.
// Generally this should not be used. Users of this package should save
// all exported fields from Client and re-use those at a later request.
// json.Marshal() and json.Unmarshal() can be used.
func SessionID(sessionID string) Option {
	return func(c *Client) {
		c.SessionID = sessionID
	}
}

// LoadBalancerID will set a load balancer ID. This is used by Garmin load
// balancers to route subsequent requests to the same backend server.
func LoadBalancerID(loadBalancerID string) Option {
	return func(c *Client) {
		c.LoadBalancerID = loadBalancerID
	}
}

// Credentials can be used to pass login credentials to NewClient.
func Credentials(email string, password string) Option {
	return func(c *Client) {
		c.Email = email
		c.Password = password
	}
}

// AutoRenewSession will set if the session should be autorenewed upon expire.
// Default is true.
func AutoRenewSession(autoRenew bool) Option {
	return func(c *Client) {
		c.autoRenewSession = autoRenew
	}
}

// DebugLogger is used to set a debug logger.
func DebugLogger(logger Logger) Option {
	return func(c *Client) {
		c.debugLogger = logger
	}
}

// DumpWriter will instruct Client to dump all HTTP requests and responses to
// and from Garmin to w.
func DumpWriter(w io.Writer) Option {
	return func(c *Client) {
		c.dumpWriter = w
	}
}

// NewClient returns a new client for accessing the unofficial Garmin Connect
// API.
func NewClient(options ...Option) *Client {
	client := &Client{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},

			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// To avoid a Cloudflare error, we have to use TLS 1.1 or 1.2.
					MinVersion: tls.VersionTLS11,
					MaxVersion: tls.VersionTLS12,
				},
			},
		},
		autoRenewSession: true,
		debugLogger:      &discardLog{},
		dumpWriter:       nil,
	}

	client.SetOptions(options...)

	return client
}

// SetOptions can be used to set various options on Client.
func (c *Client) SetOptions(options ...Option) {
	for _, option := range options {
		option(c)
	}
}

func (c *Client) dump(reqResp interface{}) {
	if c.dumpWriter == nil {
		return
	}

	var dump []byte
	switch obj := reqResp.(type) {
	case *http.Request:
		_, _ = c.dumpWriter.Write([]byte("\n\nREQUEST\n"))
		dump, _ = httputil.DumpRequestOut(obj, true)
	case *http.Response:
		_, _ = c.dumpWriter.Write([]byte("\n\nRESPONSE\n"))
		dump, _ = httputil.DumpResponse(obj, true)
	default:
		panic("unsupported type")
	}

	_, _ = c.dumpWriter.Write(dump)
}

// addCookies adds needed cookies to a http request if the values are known.
func (c *Client) addCookies(req *http.Request) {
	if c.SessionID != "" {
		req.AddCookie(&http.Cookie{
			Value: c.SessionID,
			Name:  sessionCookieName,
		})
	}

	if c.LoadBalancerID != "" {
		req.AddCookie(&http.Cookie{
			Value: c.LoadBalancerID,
			Name:  cflbCookieName,
		})
	}
}

func (c *Client) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Play nice and give Garmin engineers a way to contact us.
	req.Header.Set("User-Agent", "github.com/abrander/garmin-connect")

	// Yep. This is needed for requests sent to the API. No idea what it does.
	req.Header.Add("nk", "NT")

	c.addCookies(req)

	return req, nil
}

func (c *Client) getJSON(url string, target interface{}) error {
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	return decoder.Decode(target)
}

// write is suited for writing stuff to the API when you're NOT expected any
// data in return but a HTTP status code.
func (c *Client) write(method string, url string, payload interface{}, expectedStatus int) error {
	var body io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		body = bytes.NewReader(b)
	}

	req, err := c.newRequest(method, url, body)
	if err != nil {
		return err
	}

	// If we have a payload it is by definition JSON.
	if payload != nil {
		req.Header.Add("content-type", "application/json")
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if expectedStatus > 0 && resp.StatusCode != expectedStatus {
		return fmt.Errorf("HTTP %s returned %d (%d expected)", method, resp.StatusCode, expectedStatus)
	}

	return nil
}

// handleForbidden will try to extract an error message from the response.
func (c *Client) handleForbidden(resp *http.Response) error {
	defer resp.Body.Close()

	type proxy struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	decoder := json.NewDecoder(resp.Body)

	var errorMessage proxy

	err := decoder.Decode(&errorMessage)
	if err == nil && errorMessage.Message != "" {
		return Error(errorMessage.Message)
	}

	return ErrForbidden
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	c.debugLogger.Printf("Requesting %s at %s", req.Method, req.URL.String())

	// Save the body in case we need to replay the request.
	var save io.ReadCloser
	var err error
	if req.Body != nil {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return nil, err
		}
	}

	c.dump(req)
	t0 := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	c.dump(resp)

	// This is exciting. If the user does not have permission to access a
	// ressource, the API will return an ApplicationException and return a
	// 403 status code.
	// If the session is invalid, the Garmin API will return the same exception
	// and status code (!).
	// To distinguish between these two error cases, we look for a new session
	// cookie in the response. If a new session cookies is set by Garmin, we
	// assume our current session is invalid.
	for _, cookie := range resp.Cookies() {
		if cookie.Name == sessionCookieName {
			resp.Body.Close()
			c.debugLogger.Printf("Session invalid, requesting new session")

			// Wups. Our session got invalidated.
			c.SetOptions(SessionID(""))
			c.SetOptions(LoadBalancerID(""))

			// Re-new session.
			err = c.Authenticate()
			if err != nil {
				return nil, err
			}

			c.debugLogger.Printf("Successfully authenticated as %s", c.Email)

			// Replace the drained body
			req.Body = save

			// Replace the cookie ned newRequest with the new sessionid and load balancer key.
			req.Header.Del("Cookie")
			c.addCookies(req)

			c.debugLogger.Printf("Replaying %s request to %s", req.Method, req.URL.String())

			c.dump(req)

			// Replay the original request only once, if we fail twice
			// something is rotten, and we should give up.
			t0 = time.Now()
			resp, err = c.client.Do(req)
			if err != nil {
				return nil, err
			}
			c.dump(resp)
		}
	}

	c.debugLogger.Printf("Got HTTP status code %d in %s", resp.StatusCode, time.Since(t0).String())

	switch resp.StatusCode {
	case http.StatusBadRequest:
		resp.Body.Close()
		return nil, ErrBadRequest
	case http.StatusForbidden:
		return nil, c.handleForbidden(resp)
	case http.StatusNotFound:
		resp.Body.Close()
		return nil, ErrNotFound
	}

	return resp, err
}

// Download will retrieve a file from url using Garmin Connect credentials.
// It's mostly useful when developing new features or debugging existing
// ones.
// Please note that this will pass the Garmin session cookie to the URL
// provided. Only use this for endpoints on garmin.com.
func (c *Client) Download(url string, w io.Writer) error {
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) authenticated() bool {
	return c.SessionID != ""
}

// Authenticate using a Garmin Connect username and password provided by
// the Credentials option function.
func (c *Client) Authenticate() error {
	// We cannot use Client.do() in this function, since this function can be
	// called from do() upon session renewal.
	URL := "https://sso.garmin.com/sso/signin" +
		"?service=https%3A%2F%2Fconnect.garmin.com%2Fmodern%2F" +
		"&gauthHost=https%3A%2F%2Fconnect.garmin.com%2Fmodern%2F" +
		"&generateExtraServiceTicket=true" +
		"&generateTwoExtraServiceTickets=true"

	if c.Email == "" || c.Password == "" {
		return ErrNoCredentials
	}

	c.debugLogger.Printf("Getting CSRF token at %s", URL)

	// Start by getting CSRF token.
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	c.dump(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	c.dump(resp)

	csrfToken, err := extractCSRFToken(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	c.debugLogger.Printf("Got CSRF token: '%s'", csrfToken)

	c.debugLogger.Printf("Trying credentials at %s", URL)

	formValues := url.Values{
		"username": {c.Email},
		"password": {c.Password},
		"embed":    {"false"},
		"_csrf":    {csrfToken},
	}

	req, err = c.newRequest("POST", URL, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", URL)

	c.dump(req)

	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	c.dump(resp)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()

		return fmt.Errorf("Garmin SSO returned \"%s\"", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// Extract ticket URL
	t := regexp.MustCompile(`https:\\\/\\\/connect.garmin.com\\\/modern\\\/\?ticket=(([a-zA-Z0-9]|-)*)`)
	ticketURL := t.FindString(string(body))

	// undo escaping
	ticketURL = strings.Replace(ticketURL, "\\/", "/", -1)

	if ticketURL == "" {
		return ErrWrongCredentials
	}

	c.debugLogger.Printf("Requesting session at ticket URL %s", ticketURL)

	// Use ticket to request session.
	req, _ = c.newRequest("GET", ticketURL, nil)
	c.dump(req)
	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	c.dump(resp)
	resp.Body.Close()

	// Look for the needed sessionid cookie.
	for _, cookie := range resp.Cookies() {
		if cookie.Name == cflbCookieName {
			c.debugLogger.Printf("Found load balancer cookie with value %s", cookie.Value)

			c.SetOptions(LoadBalancerID(cookie.Value))
		}

		if cookie.Name == sessionCookieName {
			c.debugLogger.Printf("Found session cookie with value %s", cookie.Value)

			c.SetOptions(SessionID(cookie.Value))
		}
	}

	if c.SessionID == "" {
		c.debugLogger.Printf("No sessionid found")

		return ErrWrongCredentials
	}

	// The session id will not be valid until we redeem the sessions by
	// following the redirect.
	location := resp.Header.Get("Location")
	c.debugLogger.Printf("Redeeming session id at %s", location)

	req, _ = c.newRequest("GET", location, nil)
	c.dump(req)
	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	c.dump(resp)

	c.Profile, err = extractSocialProfile(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

// extractSocialProfile will try to extract the social profile from the HTML.
// This is very fragile.
func extractSocialProfile(body io.Reader) (*SocialProfile, error) {
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VIEWER_SOCIAL_PROFILE") {
			line = strings.TrimSpace(line)
			line = strings.Replace(line, "\\", "", -1)
			line = strings.TrimPrefix(line, "window.VIEWER_SOCIAL_PROFILE = ")
			line = strings.TrimSuffix(line, ";")

			profile := new(SocialProfile)

			err := json.Unmarshal([]byte(line), profile)
			if err != nil {
				return nil, err
			}

			return profile, nil
		}
	}

	return nil, errors.New("social profile not found in HTML")
}

// extractCSRFToken will try to extract the CSRF token from the signin form.
// This is very fragile. Maybe we should replace this madness by a real HTML
// parser some day.
func extractCSRFToken(body io.Reader) (string, error) {
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "name=\"_csrf\"") {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, `<input type="hidden" name="_csrf" value="`)
			line = strings.TrimSuffix(line, `" />`)

			return line, nil
		}
	}

	return "", errors.New("CSRF token not found")
}

// Signout will end the session with Garmin. If you use this for regular
// automated tasks, it would be nice to signout each time to avoid filling
// Garmin's session tables with a lot of short-lived sessions.
func (c *Client) Signout() error {
	if !c.authenticated() {
		return ErrNotAuthenticated
	}

	req, err := c.newRequest("GET", "https://connect.garmin.com/modern/auth/logout", nil)
	if err != nil {
		return err
	}

	if c.SessionID == "" {
		return nil
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.SetOptions(SessionID(""))
	c.SetOptions(LoadBalancerID(""))

	return nil
}
