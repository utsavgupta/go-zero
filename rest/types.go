package rest

import (
	"net/http"
	"net/url"
	"time"
)

type (
	// Middleware defines the middleware method.
	Middleware func(next http.HandlerFunc) http.HandlerFunc

	// A Route is a http route.
	Route struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}

	// RouteOption defines the method to customize a featured route.
	RouteOption func(r *featuredRoutes)

	jwtSetting struct {
		enabled    bool
		secret     string
		prevSecret string
	}

	oidcSetting struct {
		enabled                  bool
		client                   *http.Client
		configEndpoint           *url.URL
		introspectionEndpointKey string
		clientId                 string
		clientSecret             string
	}

	signatureSetting struct {
		SignatureConf
		enabled bool
	}

	featuredRoutes struct {
		timeout   time.Duration
		priority  bool
		jwt       jwtSetting
		oidc      oidcSetting
		signature signatureSetting
		routes    []Route
		maxBytes  int64
	}
)
