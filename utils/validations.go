package utils

import (
	"errors"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func IsNumeric(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

func ValidateMail(email string) error {

	// to do validar que sea correo

	return nil
}

func ValidatePhone(phone string) error {
	if !IsNumeric(phone) {
		return errors.New("invalid.")
	}

	if len(phone) != 9 {
		return errors.New("length invalid.")
	}

	return nil
}

func ValidateString(s string) error {
	if len(s) > 15 || len(s) == 0 {
		return errors.New("string length invalid.")
	}

	if IsNumeric(s) {
		return errors.New("string invalid")
	}

	return nil
}

func CapitalizateText(s string) string {
	loweredString := strings.ToLower(s)
	capitalizedString := cases.Title(language.Spanish).String(loweredString)

	return capitalizedString
}
