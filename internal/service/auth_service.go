package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kais-blkc/go-blog/internal/model"
	"github.com/kais-blkc/go-blog/internal/repository"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

const (
	ErrUserAlreadyExistsEmail = "пользователь с таким email уже существует"
	ErrUserAlreadyExistsName  = "пользователь с таким username уже существует"
	ErrUserWrongEmailOrPass   = "неверный email или пароль"
	ErrUserNotFound           = "пользователь не найден"
)

const (
	ClaimUserID = "user_id"
	ClaimExp    = "exp"
)

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if the user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New(ErrUserAlreadyExistsEmail)
	}

	existingUser, _ = s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New(ErrUserAlreadyExistsEmail)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Find user
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New(ErrUserWrongEmailOrPass)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New(ErrUserWrongEmailOrPass)
	}

	// Generate token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		ClaimUserID: userID,
		ClaimExp:    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	// Parser token
	token, err := jwt.Parse(tokenString, s.getKeyFunc())
	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	// Check claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	// get userID
	userID, ok := claims[ClaimUserID].(float64)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return uint(userID), nil
}

func (s *AuthService) getKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		// Check the signing method
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.jwtSecret), nil
	}
}
