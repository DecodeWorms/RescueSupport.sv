package encrypt

import "golang.org/x/crypto/bcrypt"

type Encryptor interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashed, password string) (bool, error)
}

type PasswordEncryptor struct{}

func NewPasswordEncryptor() Encryptor {
	return PasswordEncryptor{}
}

func (p PasswordEncryptor) HashPassword(password string) (string, error) {
	pas, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pas), err
}

func (p PasswordEncryptor) CompareHashAndPassword(hash, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}
