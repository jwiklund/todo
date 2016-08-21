package jira

import (
	"encoding/base64"
	"net/http"
)

func newAuth(user, pass string) http.RoundTripper {
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
	return &authTransport{basicAuth}
}

type authTransport struct {
	auth string
}

func (a *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header = map[string][]string{
		"Authorization": []string{a.auth},
	}
	return http.DefaultTransport.RoundTrip(r)
}
