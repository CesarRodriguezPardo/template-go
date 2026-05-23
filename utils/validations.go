package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ValidateMail(email string) error {
	if len(email) > 254 {
		return errors.New("invalid mail")
	}

	splittedMail := strings.Split(email, "@")
	if len(splittedMail) != 2 {
		return errors.New("invalid mail")
	}

	emailUser, domain := splittedMail[0], strings.ToLower(splittedMail[1])

	if len(emailUser) == 0 || len(emailUser) > 64 {
		return errors.New("invalid mail")
	}
	if strings.HasPrefix(emailUser, ".") || strings.HasSuffix(emailUser, ".") {
		return errors.New("invalid mail")
	}
	if strings.Contains(emailUser, "..") {
		return errors.New("invalid mail")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._+\-%]+$`).MatchString(emailUser) {
		return errors.New("invalid mail")
	}
	if len(domain) == 0 {
		return errors.New("invalid mail")
	}

	return nil
}

func IsNumeric(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

func ValidatePhone(phone string) error {
	if !IsNumeric(phone) {
		return errors.New("invalid phone number")
	}

	if len(phone) != 9 {
		return errors.New("invalid length invalid.")
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

func CapitalizeText(s string) string {
	loweredString := strings.ToLower(s)
	capitalizedString := cases.Title(language.Spanish).String(loweredString)

	return capitalizedString
}
