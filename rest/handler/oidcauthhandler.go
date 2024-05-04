package handler

import (
	"net/http"
	"net/url"
	"strings"
)

type (
	OidcAuthnConfig struct {
		OpenIDConfigurationURL url.URL
		IntrospectEndpointKey  string
		ClientId               string
		ClientSecret           string
	}
)

func OidcAuthenticate(client *http.Client, config *OidcAuthnConfig) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}

func buildRequest(config *OidcAuthnConfig) *http.Request {

	http.NewRequest(http.MethodPost, config.OpenIDConfigurationURL.RequestURI())
}

func extractAccessToken(r *http.Request) string {

	authzHeader := r.Header.Get("Authorization")

	if authzHeader == "" {
		return ""
	}

	authzHeaderTokens := strings.Split(authzHeader, " ")

	if len(authzHeaderTokens) != 2 || strings.ToLower(authzHeaderTokens[0]) != "bearer" {
		return ""
	}

	return authzHeaderTokens[1]
}
