package services

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"shopsphere/auth-service/models"
	"shopsphere/auth-service/repository"
)

var (
	ErrEmailExists = errors.New("email already registered")
)

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

func (s *AuthService) Register(
	ctx context.Context,
	req models.RegisterRequest,
) (*models.User, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrEmailExists
	}

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