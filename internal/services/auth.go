package services

import (
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/utils"
	"context"
	"errors"
)

func AuthenticateUser(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := userRepo.GetAuthDataByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if utils.CompareHashToPassword(password, user.Password) {
		return user, nil
	}

	return nil, errors.New("invalid credentials.")
}
