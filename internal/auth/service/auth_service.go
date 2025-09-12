package service

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omatheuscaetano/planus-api/internal/auth/dto"
	"github.com/omatheuscaetano/planus-api/internal/auth/model"
	"github.com/omatheuscaetano/planus-api/internal/auth/store"
	personStore "github.com/omatheuscaetano/planus-api/internal/person/store"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    store       store.AuthStore
    personStore personStore.PersonStore
}

func NewAuthService(store store.AuthStore, personStore personStore.PersonStore) *AuthService {
    return &AuthService{store: store, personStore: personStore}
}

func (s *AuthService) validateCredentials(user *model.User, password string) *errs.Error {
    if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
        return errs.InvalidCredentials()
    }

    return nil
}

func (s *AuthService) generateJWT(user *model.User) (*dto.LoginData, *errs.Error) {
    loginData := &dto.LoginData{
        ExpiresIn: time.Now().Add(24 * time.Hour * 7).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": loginData.ExpiresIn,
    })

    tokenString, err := token.SignedString([]byte(app.JWTSecret()))
    if err != nil {
        return nil, errs.From(err)
    }

    loginData.Token = tokenString
    return loginData, nil
}

func (s *AuthService) Login(c context.Context, dto *dto.Login) (*dto.LoginData, *errs.Error) {
    user, _ := s.store.FindUserByEmail(c, strings.ToLower(strings.TrimSpace(dto.Email)))

    if s.validateCredentials(user, dto.Password) != nil {
        return nil, errs.InvalidCredentials()
    }

    data, err := s.generateJWT(user)
    if err != nil {
        return nil, err
    }

    return data, nil
}

func (s *AuthService) Create(c context.Context, dto *dto.CreateUser) (*model.User, *errs.Error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errs.From(err)
    }

    return s.store.CreateUser(c, &model.User{
        Email:    strings.ToLower(strings.TrimSpace(dto.Email)),
        Password: string(hashedPassword),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    })
}
