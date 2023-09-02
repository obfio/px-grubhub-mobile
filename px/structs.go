package px

// Config is the response from the px-conf request
type Config struct {
	Enabled bool `json:"enabled"`
	IPv6    bool `json:"ipv6"`
}

// Instructions are instructions returned from px telling the client what to do, for example APPC gives you what you need for the second request
type Instructions struct {
	Do []string `json:"do"`
}

// GrubHubAuth is the response from /auth/anon when trying to get the Authorization header for the login endpoint
type GrubHubAuth struct {
	Session struct {
		AuthToken string `json:"access_token"`
	} `json:"session_handle"`
}

// LoginDebug holds info used to debug if our cookies are working on the login endpoint or not
type LoginDebug struct {
	StatusCode int      `json:"statusCode"`
	Body       string   `json:"body"`
	CookieUsed string   `json:"cookieUsed"`
	Error      string   `json:"error"`
	Payloads   []string `json:"payloads"`
}
