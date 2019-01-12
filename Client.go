package connect

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	// ErrForbidden will be returned if the client doesn't have access to the
	// requested ressource.
	ErrForbidden = errors.New("Forbidden")

	// ErrNotFound will be returned if the requested ressource could not be
	// found.
	ErrNotFound = errors.New("Not found")

	// ErrNotAuthenticated will be returned is the client is not
	// authenticated as required by the request. Remember to call
	// Authenticate().
	ErrNotAuthenticated = errors.New("Client is not authenticated")

	// ErrWrongCredentials will be returned if the username and/or
	// password is not recognized by Garmin Connect.
	ErrWrongCredentials = errors.New("Username and/or password not recognized")
)

// Client can be used to access the unofficial Garmin Connect API.
type Client struct {
	client    *http.Client
	sessionid *http.Cookie
}

// NewClient returns a new client for accessing the unofficial Garmin Connect
// API.
func NewClient() *Client {
	return &Client{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
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
	resp, err := c.client.Do(req)
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

// Authenticate using a Garmin Connect username and password.
func (c *Client) Authenticate(username string, password string) error {
	URL := "https://sso.garmin.com/sso/signin?service=https%3A%2F%2Fconnect.garmin.com%2Fmodern%2F"

	// Get ticket from Garmin SSO.
	resp, err := c.client.PostForm(URL, url.Values{
		"username": {username},
		"password": {password},
		"embed":    {"false"},
	})
	if err != nil {
		return err
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

	// Use ticket to request session.
	req, _ := c.newRequest("GET", ticketURL, nil)
	resp, err = c.do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	// Look for the needed sessionid cookie.
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SESSIONID" {
			c.sessionid = cookie
		}
	}

	if c.sessionid == nil {
		return ErrWrongCredentials
	}

	// The session id will not be valid until we redeem the sessions by
	// following the redirect.
	location := resp.Header.Get("Location")
	_, err = c.getString(location)
	if err != nil {
		return err
	}

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
