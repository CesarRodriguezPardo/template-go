package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost  = bcrypt.MinCost
	TestCost = (MinCost + MaxCost) / 2
	MaxCost  = bcrypt.MaxCost
)

func GenerateHashedPassword(password string) (string, error) {
	slicePassword := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(slicePassword, MinCost)

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
