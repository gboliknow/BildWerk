package utility

import (
	"crypto/rand"
	"math/big"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func IsValidURL(link string) bool {
	parsedURL, err := url.ParseRequestURI(link)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateOTP() (string, error) {
	otp := make([]byte, 6)
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp[i] = byte(num.Int64() + '0')
	}
	return string(otp), nil
}
