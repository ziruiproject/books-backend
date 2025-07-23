package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(password string) (string, error)
	Check(password string, hash string) bool
}

type bcryptHasher struct {
}

func NewBcryptHasher() Hasher {
	return &bcryptHasher{}
}

func (b bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b bcryptHasher) Check(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
