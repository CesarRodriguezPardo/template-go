package utils

import (
	"CesarRodriguezPardo/template-go/internal/models"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// validar mail
// validar telefono

// insertar texto con solo primera mayuscula nombre y apellido

// validar objeto completo user

func isNumeric(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

func validateMail(email string) error {

	// to do validar que sea correo

	return nil
}

func validatePhone(phone string) error {
	if len(phone) != 9 {
		return errors.New("length invalid.")
	}

	if !isNumeric(phone) {
		return errors.New("invalid.")
	}

	return nil
}

func validateString(s string) error {
	if len(s) > 15 || len(s) == 0 {
		return errors.New("string length invalid.")
	}

	if isNumeric(s) {
		return errors.New("string invalid")
	}

	return nil
}

func capitalizateText(s string) string {
	loweredString := strings.ToLower(s)
	capitalizedString := cases.Title(language.Spanish).String(loweredString)

	return capitalizedString
}

func GetValidatedUser(user *models.User) (*models.User, error) {
	if err := validateMail(user.Email); err != nil {
		return nil, fmt.Errorf("error validating mail: %w", err)
	}

	if err := validatePhone(user.Email); err != nil {
		return nil, fmt.Errorf("error validating phone number: %w", err)
	}

	if err := validateString(user.Name); err != nil {
		return nil, fmt.Errorf("error validating name: %w", err)
	}

	if err := validateString(user.MiddleName); err != nil {
		return nil, fmt.Errorf("error validating middle name: %w", err)
	}

	capitalizedName := capitalizateText(user.Name)
	capitalizedMiddleName := capitalizateText(user.MiddleName)

	user.Name = capitalizedName
	user.MiddleName = capitalizedMiddleName

	return user, nil
}
