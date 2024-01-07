package bindings

import (
	"errors"
	"net/http"
	"time"

	"github.com/Zentech-Development/conductor-proxy/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	Groups []string `json:"groups"`
	jwt.RegisteredClaims
}

func getAccessToken(username string, groups []string, expiration int) (string, error) {
	key := config.GetConfig().AccessTokenSecret

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
		return []byte(config.GetConfig().AccessTokenSecret), nil
	}, jwt.WithIssuer("conductor-proxy"))

	if err != nil {
		return AuthClaims{}, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok {
		return *claims, nil
	}

	return AuthClaims{}, errors.New("bad claims")
}

func requireAccessToken(c *gin.Context) {
	token := c.GetHeader(tokenHeaderName)
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "No access token provided",
			"data":       map[string]any{},
		})
		return
	}

	claims, err := verifyAccessToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "Bad access token provided",
			"data":       map[string]any{},
		})
		return
	}

	c.Set("userGroups", claims.Groups)
	c.Set("username", claims.Subject)
}

const (
	tokenHeaderName = "X-CONDUCTOR-KEY"
)
