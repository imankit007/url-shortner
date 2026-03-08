package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type TokenClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	UserID    string `json:"user_id"`
	TenantID  string `json:"tenant_id"`
	Email     string `json:"email"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func GenerateRS256Token(privateKey *rsa.PrivateKey, claims TokenClaims) (string, error) {
	headerPayload := struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}{
		Alg: "RS256",
		Typ: "JWT",
	}

	headerSegment, err := marshalTokenSegment(headerPayload)
	if err != nil {
		return "", err
	}

	claimsSegment, err := marshalTokenSegment(claims)
	if err != nil {
		return "", err
	}

	signingInput := headerSegment + "." + claimsSegment
	hash := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	signatureSegment := base64.RawURLEncoding.EncodeToString(signature)
	return signingInput + "." + signatureSegment, nil
}

func VerifyRS256Token(token string, publicKey *rsa.PublicKey, expectedIssuer string) (TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return TokenClaims{}, ErrInvalidTokenFormat
	}

	var header struct {
		Alg string `json:"alg"`
	}
	if err := unmarshalTokenSegment(parts[0], &header); err != nil {
		return TokenClaims{}, err
	}
	if header.Alg != "RS256" {
		return TokenClaims{}, ErrUnsupportedSigningAlgo
	}

	signingInput := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return TokenClaims{}, err
	}

	hash := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature); err != nil {
		return TokenClaims{}, err
	}

	var claims TokenClaims
	if err := unmarshalTokenSegment(parts[1], &claims); err != nil {
		return TokenClaims{}, err
	}

	return claims, nil
}

func marshalTokenSegment(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

func unmarshalTokenSegment(segment string, target any) error {
	decoded, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}

	return json.Unmarshal(decoded, target)
}
