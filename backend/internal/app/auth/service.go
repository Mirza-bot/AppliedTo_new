package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"appliedTo/internal/app/user"
	"appliedTo/internal/platform/security/password"
	"appliedTo/internal/platform/security/token"
	"appliedTo/internal/platform/validate"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service struct {
	db     *gorm.DB
	hasher password.Hasher
	jwt    *token.JWT
	users  *user.Service
}

func NewService(db *gorm.DB, hasher password.Hasher, jwt *token.JWT, users *user.Service) *Service {
	return &Service{db: db, hasher: hasher, jwt: jwt, users: users}
}

func (s *Service) Authenticate(ctx context.Context, email, plain string) (string, error) {
	normalizedEmail, err := validate.NormalizeAndValidateEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	var u user.User
	if err := s.db.WithContext(ctx).Where("email = ?", normalizedEmail).First(&u).Error; err != nil {
		return "", ErrInvalidCredentials
	}
	if !s.hasher.Verify(u.Password, plain) {
		return "", ErrInvalidCredentials
	}

	claims := map[string]any{
		"sub": strconv.FormatUint(uint64(u.ID), 10),
		"eml": u.Email,
	}
	tok, err := s.jwt.Sign(claims)
	if err != nil {
		return "", err
	}
	return tok, nil
}

func (s *Service) Register(ctx context.Context, in RegisterRequest) (string, error) {
	dto := user.UserCreateDto{
		BaseUserDto: user.BaseUserDto{
			FirstName: strings.TrimSpace(in.FirstName),
			LastName:  strings.TrimSpace(in.LastName),
			Email:     in.Email,
		},
		Password:  in.Password,
	}
	userPublic, id, err := s.users.Create(ctx, dto)
	if err != nil {
		return "", err
	}
	//auto-login on registration
	claims := map[string]any{
		"sub": strconv.FormatUint(uint64(id), 10),
		"eml": userPublic.Email,
	}
	tok, err := s.jwt.Sign(claims)
	if err != nil {
		return "", err
	}
	return tok, nil
}

func (s *Service) JWT() *token.JWT { return s.jwt }
