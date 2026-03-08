package infrastructure

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

const (
	defaultAuthServicePrivateKeyPath = "../../ops/keys/auth-service-private.pem"
	defaultAuthServicePublicKeyPath  = "../../ops/keys/auth-service-public.pem"
)

func NewAuthServicePrivateKey() (*rsa.PrivateKey, error) {
	keyPath := os.Getenv("AUTH_PRIVATE_KEY_PATH")
	if keyPath == "" {
		keyPath = defaultAuthServicePrivateKeyPath
	}

	privateKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("failed to decode auth service private key pem")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func NewAuthServicePublicKeyPEM() ([]byte, error) {
	keyPath := os.Getenv("AUTH_PUBLIC_KEY_PATH")
	if keyPath == "" {
		keyPath = defaultAuthServicePublicKeyPath
	}

	return os.ReadFile(keyPath)
}
