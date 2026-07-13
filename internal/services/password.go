package services

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordService struct{}

func NewBcryptPasswordService() *BcryptPasswordService {
	return &BcryptPasswordService{}
}

func (s *BcryptPasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *BcryptPasswordService) Verify(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
