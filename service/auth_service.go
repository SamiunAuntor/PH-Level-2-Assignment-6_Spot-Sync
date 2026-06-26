package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	apperror "spotsync/apperror"
	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
	jwtSecret      []byte
	jwtExpiresIn   time.Duration
}

type AuthClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepository repository.UserRepository, jwtSecret string, jwtExpiresIn time.Duration) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtSecret:      []byte(jwtSecret),
		jwtExpiresIn:   jwtExpiresIn,
	}
}

func (s *authService) Register(ctx context.Context, request dto.RegisterRequest) (*dto.UserResponse, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return nil, apperror.BadRequest("Validation failed", map[string]string{
			"name": "Name is required",
		}, nil)
	}

	email := normalizeEmail(request.Email)
	role := normalizeRole(request.Role)
	if role == "" {
		role = models.RoleDriver
	}

	if !isValidRole(role) {
		return nil, apperror.BadRequest("Validation failed", map[string]string{
			"role": "Role must be either driver or admin",
		}, nil)
	}

	if _, err := s.userRepository.FindByEmail(ctx, email); err == nil {
		return nil, apperror.BadRequest("Email already exists", map[string]string{
			"email": "Email already exists",
		}, nil)
	} else if !isNotFoundError(err) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return nil, apperror.Internal("Internal server error", err)
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *authService) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	email := normalizeEmail(request.Email)

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if isNotFoundError(err) {
			return nil, invalidCredentialsError(err)
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, invalidCredentialsError(err)
	}

	token, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  toLoginUserResponse(user),
	}, nil
}

func (s *authService) generateToken(userID uint, role string) (string, error) {
	now := time.Now().UTC()
	claims := AuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", apperror.Internal("Internal server error", err)
	}

	return signedToken, nil
}

func toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toLoginUserResponse(user *models.User) dto.LoginUserResponse {
	return dto.LoginUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeRole(role string) string {
	return strings.TrimSpace(role)
}

func isValidRole(role string) bool {
	return role == models.RoleDriver || role == models.RoleAdmin
}

func invalidCredentialsError(err error) error {
	return apperror.Unauthorized("Invalid email or password", nil, err)
}

func isNotFoundError(err error) bool {
	var appErr *apperror.AppError
	return errors.As(err, &appErr) && appErr.StatusCode == 404
}
