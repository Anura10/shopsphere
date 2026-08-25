package services

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"shopsphere/auth-service/models"
	"shopsphere/auth-service/repository"
	"shopsphere/auth-service/utils"
)

var ErrEmailExists = errors.New("email already registered")

type AuthService struct {
	userRepository *repository.UserRepository
}

func NewAuthService(
	userRepository *repository.UserRepository,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

// Register creates a new user.
func (s *AuthService) Register(
	ctx context.Context,
	req models.RegisterRequest,
) (*models.User, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check whether email already exists.
	existingUser, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrEmailExists
	}

	// Hash password using bcrypt.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		RoleID:       1,
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Email:        email,
		PasswordHash: string(passwordHash),
		Phone:        strings.TrimSpace(req.Phone),
		IsActive:     true,
	}

	return s.userRepository.Create(ctx, user)
}

// Login verifies credentials and generates a JWT.
func (s *AuthService) Login(
	ctx context.Context,
	req models.LoginRequest,
) (*models.User, string, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, "", errors.New("user account is inactive")
	}

	// Compare the submitted password with the stored bcrypt hash.
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT.
	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.RoleID,
	)

	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}
