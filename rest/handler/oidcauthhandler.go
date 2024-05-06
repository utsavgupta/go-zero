package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
)

type claims struct {
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}

const (
	ctxEmailKey = "email"
)

func OidcAuthenticate(client *http.Client, providerUrl *url.URL, clientId string) func(http.Handler) http.Handler {

	ctx := oidc.ClientContext(context.Background(), client)

	provider, err := oidc.NewProvider(ctx, fmt.Sprintf("%s://%s", providerUrl.Scheme, providerUrl.Host))

	if err != nil {
		panic(err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientId})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			bearerToken := extractBearerToken(r)

			if bearerToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			idToken, err := verifier.Verify(r.Context(), bearerToken)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userClaims := claims{}

			if err := idToken.Claims(&userClaims); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			fmt.Println(userClaims)

			newCtx := context.WithValue(r.Context(), ctxEmailKey, userClaims.Email)

			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

func extractBearerToken(r *http.Request) string {

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
