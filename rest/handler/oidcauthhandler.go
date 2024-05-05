package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type IntrospectionResponseBody struct {
	Active bool `json:"active"`
}

var (
	errNon2xxResponse                   = errors.New("received non 200 response")
	errCouldNotParseResponse            = errors.New("could not parse response of introspect endpoint")
	errIntrospectionKeyNotFound         = errors.New("the token introspection url was not found in the response")
	errIntrospectionUrlExtractionFailed = errors.New("could not extract the token introspection url")
)

func OidcAuthenticate(client *http.Client, openIDConfigurationURL *url.URL, introspectEndpointKey string,
	clientId string, clientSecret string) func(http.Handler) http.Handler {

	introspectionEndpoint := mustGetIntrospectionEndpointFromConfigURL(client, openIDConfigurationURL, introspectEndpointKey)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			accessToken := extractAccessToken(r)

			fmt.Println(accessToken)

			if accessToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			oidcReq := buildOidcRequest(introspectionEndpoint, clientId, clientSecret, accessToken)
			fmt.Printf("%v\n", oidcReq.URL.String())
			oidcResponse, err := makeIntrospectionRequest(client, oidcReq)

			fmt.Printf("%v\n", oidcResponse)

			if err != nil || !oidcResponse.Active {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func mustGetIntrospectionEndpointFromConfigURL(client *http.Client, configUrl *url.URL, introspectionEndpointKey string) string {

	configRequest, _ := http.NewRequest(http.MethodGet, configUrl.String(), nil)

	configResponse, err := client.Do(configRequest)

	if err != nil {
		panic(err)
	}

	defer configResponse.Body.Close()

	if configResponse.StatusCode != http.StatusOK {
		panic(errNon2xxResponse)
	}

	responseBody := make(map[string]any)

	if err := json.NewDecoder(configResponse.Body).Decode(&responseBody); err != nil {
		panic(err)
	}

	introspectionEndpointVal, ok := responseBody[introspectionEndpointKey]

	if !ok {
		panic(errIntrospectionKeyNotFound)
	}

	introspectionEndpoint, ok := introspectionEndpointVal.(string)

	if !ok {
		panic(errIntrospectionUrlExtractionFailed)
	}

	return introspectionEndpoint
}

func buildOidcRequest(introspectionEndpoint, clientId, clientSecret, accessToken string) *http.Request {

	values := url.Values{}

	values.Set("token", accessToken)
	values.Set("token_type_hint", "access_token")

	r, _ := http.NewRequest(http.MethodPost, introspectionEndpoint, strings.NewReader(values.Encode()))

	r.SetBasicAuth(clientId, clientSecret)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return r
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

func makeIntrospectionRequest(client *http.Client, introspectRequest *http.Request) (*IntrospectionResponseBody, error) {

	resp, err := client.Do(introspectRequest)

	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%v\n", resp.StatusCode)
		return nil, errNon2xxResponse
	}

	defer resp.Body.Close()

	var introspectionResponseBody IntrospectionResponseBody

	if err := json.NewDecoder(resp.Body).Decode(&introspectionResponseBody); err != nil {
		fmt.Printf("%v\n", err)
		return nil, errCouldNotParseResponse
	}

	return &introspectionResponseBody, nil
}
