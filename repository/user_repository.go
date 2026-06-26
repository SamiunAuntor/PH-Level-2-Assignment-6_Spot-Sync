package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	apperror "spotsync/apperror"
	"spotsync/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateEmailError(err) {
			return apperror.BadRequest("Email already exists", map[string]string{
				"email": "Email already exists",
			}, err)
		}

		return apperror.Internal("Internal server error", err)
	}

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("User not found", nil, err)
		}

		return nil, apperror.Internal("Internal server error", err)
	}

	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("User not found", nil, err)
		}

		return nil, apperror.Internal("Internal server error", err)
	}

	return &user, nil
}

func isDuplicateEmailError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate key") && strings.Contains(message, "email")
}
