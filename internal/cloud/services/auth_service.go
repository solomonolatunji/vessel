package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"vessel.dev/vessel/internal/cloud/repos"
	"vessel.dev/vessel/internal/models"
)

// Sentinel errors for auth operations.
var (
	ErrEmailTaken         = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotVerified   = errors.New("email not verified — please check your inbox")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrInvalidToken       = errors.New("invalid or expired verification token")
)

// AuthService handles all authentication business logic.
type AuthService struct {
	repo   repos.AuthRepo
	mailer *MailerService
}

// NewAuthService creates an AuthService. mailer may be nil if SES is not configured.
func NewAuthService(repo repos.AuthRepo, mailer *MailerService) *AuthService {
	return &AuthService{repo: repo, mailer: mailer}
}

// Register creates a new cloud user, sends a welcome email, and returns a JWT.
func (s *AuthService) Register(ctx context.Context, email, password, fullName string) (string, error) {
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("checking email: %w", err)
	}
	if existing != nil {
		return "", ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}

	verifyToken, err := generateHex(32)
	if err != nil {
		return "", fmt.Errorf("generating verify token: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	user := &models.CloudUser{
		ID:                   uuid.New().String(),
		Email:                email,
		FullName:             fullName,
		PasswordHash:         string(hash),
		Role:                 "user",
		EmailVerified:        false,
		VerifyToken:          verifyToken,
		VerifyTokenExpiresAt: &expiresAt,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("creating user: %w", err)
	}

	baseURL := os.Getenv("VESSEL_CLOUD_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	verifyURL := fmt.Sprintf("%s/api/cloud/auth/verify-email?token=%s", baseURL, verifyToken)

	if s.mailer != nil {
		// Non-fatal — log internally but don't fail registration.
		if err := s.mailer.SendWelcomeEmail(ctx, email, fullName, verifyURL); err != nil {
			_ = err
		}
	}

	return s.generateJWT(user)
}

// Login authenticates a cloud user and returns a signed JWT.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("looking up user: %w", err)
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.EmailVerified {
		return "", ErrEmailNotVerified
	}

	return s.generateJWT(user)
}

// ForgotPassword generates a 6-digit OTP and sends it to the user's email.
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("looking up user: %w", err)
	}
	if user == nil {
		// Don't leak whether the email exists.
		return nil
	}

	otp, err := generateOTP()
	if err != nil {
		return fmt.Errorf("generating OTP: %w", err)
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	if err := s.repo.SaveOTP(ctx, user.ID, otp, expiresAt); err != nil {
		return fmt.Errorf("saving OTP: %w", err)
	}

	if s.mailer != nil {
		if err := s.mailer.SendOTPResetEmail(ctx, email, user.FullName, otp, "15 minutes"); err != nil {
			_ = err
		}
	}

	return nil
}

// ResetPassword validates the OTP and updates the user's password.
func (s *AuthService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("looking up user: %w", err)
	}
	if user == nil {
		return ErrInvalidOTP
	}

	if user.OTPCode != otp {
		return ErrInvalidOTP
	}
	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return ErrInvalidOTP
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	if err := s.repo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	return s.repo.ClearOTP(ctx, user.ID)
}

// VerifyEmail marks a user's email as verified using the provided token.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.repo.GetUserByVerifyToken(ctx, token)
	if err != nil {
		return fmt.Errorf("looking up token: %w", err)
	}
	if user == nil {
		return ErrInvalidToken
	}
	if user.VerifyTokenExpiresAt == nil || time.Now().After(*user.VerifyTokenExpiresAt) {
		return ErrInvalidToken
	}

	return s.repo.MarkEmailVerified(ctx, user.ID)
}

func (s *AuthService) generateJWT(user *models.CloudUser) (string, error) {
	secret := os.Getenv("VESSEL_CLOUD_JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("signing JWT: %w", err)
	}
	return signed, nil
}

func generateHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generateOTP() (string, error) {
	var otp string
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += n.String()
	}
	return otp, nil
}
