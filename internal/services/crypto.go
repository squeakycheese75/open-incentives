package services

import (
	"crypto/rand"
	"fmt"
)

type CryptoService struct {
}

func NewCryptoService() *CryptoService {
	return &CryptoService{}
}

func (CryptoService) GenerateKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}
