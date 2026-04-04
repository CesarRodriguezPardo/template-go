package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     = bcrypt.MinCost
	DefaultCost = 12
	MaxCost     = bcrypt.MaxCost
)

func GenerateHash(password string) (string, error) {
	slicePassword := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(slicePassword, DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHashToPassword(password string, hashedPassword string) bool {
	slicePassword := []byte(password)
	sliceHashedPassword := []byte(hashedPassword)
	success := bcrypt.CompareHashAndPassword(sliceHashedPassword, slicePassword)

	if success == nil {
		return true
	}
	return false
}
