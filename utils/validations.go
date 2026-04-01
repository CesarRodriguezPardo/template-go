package utils

import (
	"CesarRodriguezPardo/template-go/internal/models"
	"errors"
	"strconv"
	"strings"
)

func ToLowerEmail(email string) string {
	lowerEmail := strings.ToLower(email)
	return lowerEmail
}

// validacion super simple, deberia mejorarse pero es un alcance a lo que se espera
func GetValidatedMail(email string) (string, error) {
	splittedMail := strings.Split(email, "@")

	if len(splittedMail) != 2 || email == "" {
		return email, errors.New("Email debe seguir el formato example@usach.cl")
	}

	if splittedMail[1] != "usach.cl" {
		return email, errors.New("Email debe ser institucional Usach.")
	}

	return email, nil
}

func ValidatePhone(phone string) error {
	intPhone, err := strconv.Atoi(phone)

	if err != nil {
		return errors.New("El número de teléfono debe ser numérico.")
	}

	// se asume formato 9 1234 5678
	if len(phone) != 9 || intPhone < 0 {
		return errors.New("El número de teléfono no tiene el formato correcto.")
	}

	return nil
}

func ValidateUserObject(user *models.User) error {
	phone := user.Phone
	err := ValidatePhone(phone)
	if err != nil {
		return err
	}

	email := user.Email
	_, err = GetValidatedMail(email)
	if err != nil {
		return err
	}

	if user.Name == "" {
		return errors.New("El nombre del usuario no puede ser vacio.")
	}
	if user.MiddleName == "" {
		return errors.New("El apellido del usuario no puede ser vacio.")
	}

	return nil
}

func ValidateUserPostgresObject(user *models.User) error {
	phone := user.Phone
	err := ValidatePhone(phone)
	if err != nil {
		return err
	}

	email := user.Email
	email, err = GetValidatedMail(email)
	if err != nil {
		return err
	}

	if user.Name == "" {
		return errors.New("El nombre del usuario no puede ser vacio.")
	}
	if user.MiddleName == "" {
		return errors.New("El apellido del usuario no puede ser vacio.")
	}

	return nil
}
