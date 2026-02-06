package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the claims for Epic backend OAuth2
// Note: Epic expects 'aud' as a string, not an array
type JWTClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"` // String, not array - Epic requirement
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	JWTID     string `json:"jti"`
}

func (c JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

func (c JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

func (c JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c JWTClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c JWTClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{c.Audience}, nil
}

func GenerateJWT(clientID, tokenURL, privateKeySource string) (string, error) {
	privateKey, err := loadPrivateKey(privateKeySource)
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %w", err)
	}

	now := time.Now()
	claims := JWTClaims{
		Issuer:    clientID,
		Subject:   clientID,
		Audience:  tokenURL, // String, not array
		ExpiresAt: now.Add(5 * time.Minute).Unix(),
		IssuedAt:  now.Unix(),
		JWTID:     uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS384, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}

// loadPrivateKey loads an RSA private key from either:
// - Environment variable: if source starts with "env:" (e.g., "env:FHIR_PRIVATE_KEY")
// - File path: otherwise treats source as a file path
func loadPrivateKey(source string) (*rsa.PrivateKey, error) {
	var keyData []byte
	var err error

	if strings.HasPrefix(source, "env:") {
		// Read from environment variable
		envVar := strings.TrimPrefix(source, "env:")
		keyData = []byte(os.Getenv(envVar))
		if len(keyData) == 0 {
			return nil, fmt.Errorf("environment variable %s is not set or empty", envVar)
		}
	} else {
		// Read from file
		// Expand ~ to home directory
		if strings.HasPrefix(source, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			source = strings.Replace(source, "~", home, 1)
		}

		keyData, err = os.ReadFile(source)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block - ensure the key is in PEM format")
	}

	// Try PKCS#8 first, then PKCS#1
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("private key is not RSA")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return key, nil
}
