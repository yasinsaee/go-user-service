package jwt

import (
	"crypto/rsa"
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
		ID       string   `json:"id"`
		Username string   `json:"username"`
		Roles    []string `json:"roles"`
		Access   []string `json:"access"`
		jwt.StandardClaims
	}

	JWTConfig struct {
		PrivateKey []byte
		PublicKey  []byte
		Exp        int
	}

	TokenConfig struct {
		ID       string   `json:"id"`
		Username string   `json:"username"`
		Roles    []string `json:"roles"`
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
	expirationTime := time.Now().UTC().Add(time.Hour * time.Duration(conf.Exp))
	return t.generateToken(expirationTime, []byte(conf.PrivateKey))
}

func (t *TokenConfig) generateToken(expirationTime time.Time, privateKey []byte) (string, time.Time, error) {
	claims := &JWTClaims{
		ID:       t.ID,
		Username: t.Username,
		Roles:    t.Roles,
		Access:   t.Access,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", time.Now().UTC(), err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", time.Now().UTC(), err
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
	token = strings.TrimPrefix(token, "Bearer ")
	claims := &JWTClaims{}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(conf.PublicKey)
	if err != nil {
		return nil, err
	}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

func GetPublicKey() (*rsa.PublicKey, error) {
	return jwt.ParseRSAPublicKeyFromPEM(conf.PublicKey)
}
