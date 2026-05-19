package service

import (
	"Server/internal/domain"
	"Server/internal/dto"
	"Server/internal/model"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	authRepo domain.AuthRepository
}

func NewAuthService(repo domain.AuthRepository) *authService {
	return &authService{authRepo: repo}
}

func (s *authService) Register(c fiber.Ctx, req dto.RegisterRequest) (string, string, error) {
	existing, err := s.authRepo.FindUserByEmail(req.Email)

	if err != nil && err.Error() != "user not found" {
		return "", "", err
	}

	if existing != nil {
		return "", "", fmt.Errorf("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("error hashing password")
	}

	user := model.User{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hash),
		Status:    1,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	u, err := s.authRepo.CreateUser(user)

	if err != nil {
		return "", "", fmt.Errorf("error creating user: %w", err)
	}

	token, err := generateJWT(u)
	if err != nil {
		return "", "", fmt.Errorf("error generating token: %w", err)
	}

	return token, u.ID.Hex(), nil
}

func (s *authService) Login(c fiber.Ctx, req dto.LoginRequest) (string, string, error) {
	existing, err := s.authRepo.FindUserByEmail(req.Email)

	if err != nil {
		return "", "", fmt.Errorf("error finding user")
	}

	if existing == nil {
		return "", "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(req.Password)); err != nil {
		return "", "", fmt.Errorf("invalid password")
	}

	token, err := generateJWT(existing)
	if err != nil {
		return "", "", fmt.Errorf("error generating token: %w", err)
	}

	return token, existing.ID.Hex(), nil
}
