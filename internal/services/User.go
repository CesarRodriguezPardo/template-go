package services

import (
	"CesarRodriguezPardo/template-go/internal/dto"
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	uuid "github.com/satori/go.uuid"
)

func validateUserParams(user *models.User) error {
	if err := utils.ValidateMail(user.Email); err != nil {
		return errors.New(err.Error())
	}

	if err := utils.ValidatePhone(user.Phone); err != nil {
		return errors.New(err.Error())
	}

	if err := utils.ValidateString(user.Name); err != nil {
		return errors.New(err.Error())
	}

	if err := utils.ValidateString(user.MiddleName); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func capitalizeUserParams(user *models.User) {
	capitalizedName := utils.CapitalizeText(user.Name)
	capitalizedMiddleName := utils.CapitalizeText(user.MiddleName)

	user.Name = capitalizedName
	user.MiddleName = capitalizedMiddleName
}

func CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := req.ToModel()

	if err := validateUserParams(user); err != nil {
		return nil, fmt.Errorf("error validating user: %w", err)
	}
	capitalizeUserParams(user)

	if user.Role == "" {
		user.Role = string(models.WORKER)
	}

	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	user.Password = hashedPass

	id, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "email") {
				return nil, errors.New("user already exists with email")
			}
			if strings.Contains(pgErr.ConstraintName, "phone") {
				return nil, errors.New("user already exists with phone number")
			}
		}
		return nil, errors.New("could not create user")
	}

	user.ID = id
	return dto.UserToResponse(user), nil
}

func GetAllUsers(ctx context.Context, limit, offset int) ([]*dto.UserResponse, int, error) {
	users, total, err := userRepo.GetAllUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("could not get users: %w", err)
	}

	return dto.UsersToResponseList(users), total, nil
}

func GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	return dto.UserToResponse(user), nil
}

func UpdateUser(ctx context.Context, requesterRole string, requesterID uuid.UUID, targetUserID uuid.UUID, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check permissions
	if requesterRole != string(models.ADMIN) {
		if requesterID != targetUserID {
			return nil, errors.New("unauthorized to edit this user")
		}
	}

	existingUser, err := userRepo.GetUserByID(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	updatedUser := req.ToModel()

	if err := validateUserParams(updatedUser); err != nil {
		return nil, fmt.Errorf("error validating user: %w", err)
	}
	capitalizeUserParams(updatedUser)

	updatedUser.ID = targetUserID

	// Only admins can change roles; otherwise preserve existing role
	if requesterRole != string(models.ADMIN) {
		updatedUser.Role = existingUser.Role
	}

	err = userRepo.UpdateUser(ctx, updatedUser)
	if err != nil {
		return nil, fmt.Errorf("could not update user: %w", err)
	}

	return dto.UserToResponse(updatedUser), nil
}

func DeleteUser(ctx context.Context, requesterRole string, targetUserID uuid.UUID) error {
	// Check permissions
	if requesterRole != string(models.ADMIN) {
		return errors.New("unauthorized to delete users")
	}

	err := userRepo.DeleteUser(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	return nil
}
