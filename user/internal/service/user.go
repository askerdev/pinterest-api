package service

import (
	"context"
	"github.com/askerdev/pinterest.user/internal/domain"
	"github.com/askerdev/pinterest.user/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

type User struct {
	userRepo repository.User
	hasher   Hasher
}

func NewUser(userRepo repository.User, hasher Hasher) *User {
	return &User{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (s *User) Signup(ctx context.Context, username, password string) error {
	hashedPass, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	_, err = s.userRepo.Create(ctx, username, hashedPass)
	return err
}

func (s *User) Signin(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if !s.hasher.Compare(password, user.Password) {
		return "", domain.ErrInvalidPassword
	}

	claims := jwt.MapClaims{
		"subject": user.ID,
		"exp":     time.Now().Add(viper.GetDuration("jwt_refresh_exp")).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(viper.GetString("jwt_refresh_secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (s *User) AccessToken(ctx context.Context, id string) (string, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"subject":   user.ID,
		"username":  user.Username,
		"photo_url": user.PhotoUrl,
		"exp":       time.Now().Add(viper.GetDuration("jwt_access_exp")).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(viper.GetString("jwt_access_secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
