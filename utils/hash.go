package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// constantes relativas al costo computacional del hasheo.

const (
	MinCost  = bcrypt.MinCost
	TestCost = (MinCost + MaxCost) / 2 // unicamente para testear velocidad
	MaxCost  = bcrypt.MaxCost
)

// GeneratePassword: funcion para generar password hasheada.
func GeneratePassword(password string) string {
	slicePassword := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(slicePassword, MinCost)

	if err != nil {
		logger.Fatal("No se pudo generar el hash")
	}

	return string(hash)
}

// CompareHashToPassword: funcion para verificar si una password es la que esta hasheada.
func CompareHashToPassword(password string, hashedPassword string) bool {
	slicePassword := []byte(password)
	sliceHashedPassword := []byte(hashedPassword)
	success := bcrypt.CompareHashAndPassword(sliceHashedPassword, slicePassword)

	if success == nil {
		return true
	}
	return false
}
