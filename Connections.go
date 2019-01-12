package connect

// Connections will list the connections of displayName. If displayName is
// empty, the current authenticated connection list wil be returned.
func (c *Client) Connections(displayName string) ([]SocialProfile, error) {
	URL := "https://connect.garmin.com/modern/proxy/userprofile-service/socialProfile/connections/" + displayName

	if !c.authenticated() && displayName == "" {
		return nil, ErrNotAuthenticated
	}

	var proxy struct {
		Connections []SocialProfile `json:"userConnections"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, err
	}

	return proxy.Connections, nil
}

// PendingConnections returns a list of pending connections.
func (c *Client) PendingConnections() ([]SocialProfile, error) {
	URL := "https://connect.garmin.com/modern/proxy/userprofile-service/connection/pending"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	pending := make([]SocialProfile, 0, 10)

	err := c.getJSON(URL, &pending)
	if err != nil {
		return nil, err
	}

	return pending, nil
}
