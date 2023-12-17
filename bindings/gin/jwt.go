package bindings

import (
	"errors"
	"time"

	"github.com/Zentech-Development/conductor-proxy/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	Groups []string `json:"groups"`
	jwt.RegisteredClaims
}

func getAccessToken(username string, groups []string, expiration int) (string, error) {
	key := config.GetConfig().SecretKey

	claims := AuthClaims{
		groups,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "conductor-proxy",
			Subject:   username,
			Audience:  []string{"conductor-proxy"},
			ID:        "1",
		},
	}

	if expiration > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(expiration) * time.Second))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	return signedToken, err
}

func verifyAccessToken(signedToken string) (AuthClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return config.GetConfig().SecretKey, nil
	}, jwt.WithIssuer("conductor-proxy"), jwt.WithExpirationRequired(), jwt.WithLeeway(3*time.Second))

	if err != nil {
		return AuthClaims{}, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok {
		return *claims, nil
	}

	return AuthClaims{}, errors.New("Bad claims")
}
