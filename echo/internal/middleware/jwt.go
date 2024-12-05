package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

var JwksURL = "https://url-to-your-jwks.com"

func getKey(token *jwt.Token) (interface{}, error) {
	keySet, err := jwk.Fetch(context.Background(), JwksURL)
	if err != nil {
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have a key ID in the kid field")
	}

	key, found := keySet.LookupKeyID(keyID)
	if !found {
		return nil, fmt.Errorf("unable to find key %q", keyID)
	}

	var pubkey rsa.PublicKey
	if err := jwk.Export(key, &pubkey); err != nil {
		return nil, fmt.Errorf("unable to export the public key: %w", err)
	}

	return &pubkey, nil
}

func skipper(c echo.Context) bool {
	skipPaths := []string{"/ping", "/health", "/status"}
	for _, path := range skipPaths {
		if c.Path() == path {
			return true
		}
	}
	return false
}

var JWT = echojwt.WithConfig(echojwt.Config{
	KeyFunc:       getKey,
	SigningMethod: "RS256",
	Skipper:       skipper,
})
