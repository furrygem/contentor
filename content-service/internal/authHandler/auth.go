package authHandler

import (
	"crypto"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	EdPublicKey crypto.PublicKey
	config      *AuthHandlerConfig
}

func New(c *AuthHandlerConfig) (*AuthHandler, error) {
	ah := &AuthHandler{
		config: c,
	}
	err := ah.loadKey()
	if err != nil {
		return nil, err
	}
	return ah, nil
}

func (ah *AuthHandler) loadKey() error {
	file, err := os.ReadFile(ah.config.EdDSAPublicKeyFile)
	if err != nil {
		return err
	}
	key, err := jwt.ParseEdPublicKeyFromPEM(file)
	if err != nil {
		return err
	}
	ah.EdPublicKey = key
	return nil
}

func (ah *AuthHandler) ParseJWT(token string) (*jwt.Token, jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method == jwt.SigningMethodEdDSA {
			return ah.EdPublicKey, nil
		}
		return nil, jwt.ErrInvalidKey
	}, jwt.WithValidMethods([]string{"EdDSA"}))
	claims, _ := ah.getClaims(t)
	return t, claims, err
}

func (ah *AuthHandler) getClaims(t *jwt.Token) (jwt.MapClaims, bool) {
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, true
	}
	return nil, false
}
