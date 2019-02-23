package connect

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrForbidden will be returned if the client doesn't have access to the
	// requested ressource.
	ErrForbidden = errors.New("Forbidden")

	// ErrNotFound will be returned if the requested ressource could not be
	// found.
	ErrNotFound = errors.New("Not found")

	// ErrNoCredentials will be returned if credentials are needed - but none
	// are set.
	ErrNoCredentials = errors.New("No credentials set")

	// ErrNotAuthenticated will be returned is the client is not
	// authenticated as required by the request. Remember to call
	// Authenticate().
	ErrNotAuthenticated = errors.New("Client is not authenticated")

	// ErrWrongCredentials will be returned if the username and/or
	// password is not recognized by Garmin Connect.
	ErrWrongCredentials = errors.New("Username and/or password not recognized")
)

const (
	// sessionCookieName is the magic session cookie name.
	sessionCookieName = "SESSIONID"
)

// Client can be used to access the unofficial Garmin Connect API.
type Client struct {
	client           *http.Client
	sessionid        *http.Cookie
	login            string
	password         string
	autoRenewSession bool
	debugLogger      Logger
	dumpWriter       io.Writer
}

// SessionID will set a predefined session ID. This can be useful for clients
// keeping state. A few HTTP roundtrips can be saved, if the session ID is
// reused. And some load wouyld be taken of Garmin servers.
func SessionID(sessionID string) func(*Client) {
	return func(c *Client) {
		if sessionID != "" {
			c.sessionid = &http.Cookie{
				Value: sessionID,
				Name:  sessionCookieName,
			}
		}
	}
}

// Credentials can be used to pass login credentials to NewClient.
func Credentials(email string, password string) func(*Client) {
	return func(c *Client) {
		c.login = email
		c.password = password
	}
}

// AutoRenewSession will set if the session should be autorenewed upon expire.
// Default is true.
func AutoRenewSession(autoRenew bool) func(*Client) {
	return func(c *Client) {
		c.autoRenewSession = autoRenew
	}
}

// DebugLogger is used to set a debug logger.
func DebugLogger(logger Logger) func(*Client) {
	return func(c *Client) {
		c.debugLogger = logger
	}
}

// DumpWriter will instruct Client to dump all HTTP requests and responses to
// and from Garmin to w.
func DumpWriter(w io.Writer) func(*Client) {
	return func(c *Client) {
		c.dumpWriter = w
	}
}

// NewClient returns a new client for accessing the unofficial Garmin Connect
// API.
func NewClient(options ...func(*Client)) *Client {
	client := &Client{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
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
func (c *Client) SetOptions(options ...func(*Client)) {
	for _, option := range options {
		option(c)
	}
}

// SessionID returns the current known session ID.
func (c *Client) SessionID() string {
	if c.sessionid != nil {
		return c.sessionid.Value
	}

	return ""
}

func (c *Client) dump(reqResp interface{}) {
	if c.dumpWriter == nil {
		return
	}

	var dump []byte
	switch reqResp.(type) {
	case *http.Request:
		c.dumpWriter.Write([]byte("\n\nREQUEST\n"))
		dump, _ = httputil.DumpRequestOut(reqResp.(*http.Request), true)
	case *http.Response:
		c.dumpWriter.Write([]byte("\n\nRESPONSE\n"))
		dump, _ = httputil.DumpResponse(reqResp.(*http.Response), true)
	default:
		panic("unsupported type")
	}

	c.dumpWriter.Write(dump)
}

func (c *Client) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Play nice and give Garmin engineers a way to contact us.
	req.Header.Set("User-Agent", "github.com/abrander/garmin-connect")

	// If sessionid is known, add the cookie.
	if c.sessionid != nil {
		req.AddCookie(c.sessionid)
	}

	return req, nil
}

func (c *Client) getJSON(URL string, target interface{}) error {
	req, err := c.newRequest("GET", URL, nil)
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

func (c *Client) getString(URL string) (string, error) {
	req, err := c.newRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	c.debugLogger.Printf("Requesting %s at %s", req.Method, req.URL.String())

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
			c.debugLogger.Printf("Session invalid, requesting new session")

			// Wups. Our session got invalidated.
			c.sessionid = nil

			// Re-new session.
			err = c.Authenticate()
			if err != nil {
				return nil, err
			}

			c.debugLogger.Printf("Successfully authenticated as %s", c.login)

			// Replace the cookie ned newRequest with the new sessionid.
			req.Header.Del("Cookie")
			req.AddCookie(c.sessionid)

			c.debugLogger.Printf("Replaying %s request to %s", req.Method, req.URL.String())

			c.dump(req)

			// Replay the original request only once, if we fail twice
			// something is rotten, and we should give up.
			t0 = time.Now()
			resp, err = c.client.Do(req)
			if err != nil {
				return nil, err
			}
		}
	}

	c.dump(resp)

	c.debugLogger.Printf("Got HTTP status code %d in %s", resp.StatusCode, time.Since(t0).String())

	switch resp.StatusCode {
	case http.StatusForbidden:
		resp.Body.Close()
		return nil, ErrForbidden
	case http.StatusNotFound:
		resp.Body.Close()
		return nil, ErrNotFound
	}

	return resp, err
}

func (c *Client) authenticated() bool {
	return c.sessionid != nil
}

// Authenticate using a Garmin Connect username and password provided by
// the Credentials option function.
func (c *Client) Authenticate() error {
	// We cannot use Client.do() in this function, since this function can be
	// called from do() upon session renewal.
	URL := "https://sso.garmin.com/sso/signin?service=https%3A%2F%2Fconnect.garmin.com%2Fmodern%2F"

	if c.login == "" || c.password == "" {
		return ErrNoCredentials
	}

	c.debugLogger.Printf("Trying credentials at %s", URL)

	// Get ticket from Garmin SSO.
	resp, err := c.client.PostForm(URL, url.Values{
		"username": {c.login},
		"password": {c.password},
		"embed":    {"false"},
	})
	if err != nil {
		return err
	}
	c.dump(resp)

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
	req, _ := c.newRequest("GET", ticketURL, nil)
	c.dump(req)
	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	c.dump(resp)
	resp.Body.Close()

	// Look for the needed sessionid cookie.
	for _, cookie := range resp.Cookies() {
		if cookie.Name == sessionCookieName {
			c.debugLogger.Printf("Found session cookie with value %s", cookie.Value)

			c.sessionid = cookie
		}
	}

	if c.sessionid == nil {
		c.debugLogger.Printf("No sessionid found")

		return ErrWrongCredentials
	}

	// The session id will not be valid until we redeem the sessions by
	// following the redirect.
	location := resp.Header.Get("Location")
	c.debugLogger.Printf("Redeeming session id at %s", location)

	req, _ = c.newRequest("GET", location, nil)
	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	c.dump(resp)
	resp.Body.Close()

	return nil
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

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.sessionid = nil

	return nil
}
