package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	conf            JWTConfig
	ErrTokenExpired = errors.New("token expired")
)

type (
	// JWTClaims is a struct that will be encoded to a JWT.
	JWTClaims struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Access   string `json:"access"`
		jwt.StandardClaims
	}

	JWTConfig struct {
		Secret string
		Exp    int
	}

	TokenConfig struct {
		ID       string   `json:"id"`
		Username string   `json:"username"`
		Role     string   `json:"role"`
		Access   []string `json:"access"`
	}
)

type Model interface {
	Get(filter bson.M) error
}

func Init(config JWTConfig) {
	conf = config
}

func (t *TokenConfig) GenerateAccessToken() (string, time.Time, error) {
	// Declare the expiration time of the token - ? hours.
	expirationTime := time.Now().Add(time.Hour * time.Duration(conf.Exp))
	return t.generateToken(expirationTime, []byte(conf.Secret))
}

func (t *TokenConfig) generateToken(expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &JWTClaims{
		ID:       t.ID,
		Username: t.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expirationTime, nil
}

func CurrentToken(c *echo.Context) (*JWTClaims, error) {
	token := (*c).Request().Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("token not found")
	}
	return validation(token)
}

func validation(token string) (*JWTClaims, error) {
	// Initialize a new instance of `Claims`
	token = strings.TrimPrefix(token, "Bearer ")

	claims := &JWTClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// check expiration time
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	return claims, nil
}
