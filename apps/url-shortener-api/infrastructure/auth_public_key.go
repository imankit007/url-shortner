package infrastructure

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

const defaultAuthServicePublicKeyPath = "../../ops/keys/auth-service-public.pem"

func NewAuthServicePublicKey() (*rsa.PublicKey, error) {
	keyPath := os.Getenv("AUTH_PUBLIC_KEY_PATH")
	if keyPath == "" {
		keyPath = defaultAuthServicePublicKeyPath
	}

	publicKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, errors.New("failed to decode auth service public key pem")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("auth service public key is not rsa")
	}

	return publicKey, nil
}

func NewAuthTokenIssuer() string {
	return "url-shortener-auth-service"
}
