package service

import (
	"Server/internal/domain"
	"Server/internal/dto"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo domain.AuthRepository
}

func NewAuthService(repo domain.AuthRepository) *authService {
	return &authService{repo: repo}
}

func (s *authService) Register(c fiber.Ctx, req dto.RegisterRequest) (string, error) {
	existing, err := s.repo.FindUserByEmail(req.Email)

	if err != nil && err.Error() != "user not found" {
		return "", err
	}

	if existing != nil {
		return "", fmt.Errorf("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password")
	}

	user := domain.User{
		ID:        primitive.NewObjectID(),
		Name:      req.Name + " " + req.LastName,
		Email:     req.Email,
		Password:  string(hash),
		Status:    1,
		Followers: make([]primitive.ObjectID, 0),
		Following: make([]primitive.ObjectID, 0),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	u, err := s.repo.CreateUser(user)

	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	println(user.ID.Hex())

	return generateJWT(u)
}

func (s *authService) Login(c fiber.Ctx, req dto.LoginRequest) (string, error) {
	existing, err := s.repo.FindUserByEmail(req.Email)

	if err != nil {
		return "", fmt.Errorf("error finding user")
	}

	if existing == nil {
		return "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	return generateJWT(existing)
}
