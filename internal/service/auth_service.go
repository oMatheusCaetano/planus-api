package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omatheuscaetano/planus-api/internal/dto"
	"github.com/omatheuscaetano/planus-api/pkg/env"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
    return &AuthService{userService: userService}
}

func (s *AuthService) Login(c context.Context, props *dto.LoginData) (*dto.LoginResponse, *errs.Error) {
    user, _ := s.userService.FindByUsername(c, props.Username)

    if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(props.Password)) != nil {
        return nil, errs.InvalidCredentials()
    }

    expiration := time.Now().Add(24 * time.Hour * 7)
    loginData := &dto.LoginResponse{
        User:      user,
        ExpiresIn: expiration.Unix(),
    }

    claims := dto.JWTClaims{
        Sub: user.ID.String(),
        RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiration)},
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte(env.JWTSecret()))
    if err != nil {
        return nil, errs.From(err)
    }

    loginData.Token = tokenString
    return loginData, nil
}
