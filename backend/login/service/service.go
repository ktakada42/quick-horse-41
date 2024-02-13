package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	lr "app/login/repository"
	ur "app/user/repository"
)

type LoginService interface {
	IsPasswordCorrect(userId, password string) (bool, error)
	CreateToken(userId string) (token string, expiresIn string, err error)
	UpdateToken(userId string) (token string, expiresIn string, err error)
}

type loginService struct {
	lr lr.LoginRepository
	ur ur.UserRepository
}

func NewLoginService(lr lr.LoginRepository, ur ur.UserRepository) LoginService {
	return &loginService{lr: lr, ur: ur}
}

func (s *loginService) IsPasswordCorrect(userId, password string) (bool, error) {
	hashedPassword, err := s.ur.GetHashedPassword(userId)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func (s *loginService) CreateToken(userId string) (token string, expiresIn string, err error) {
	token, expiresIn, err = createJWT(userId)
	if err != nil {
		return "", "", err
	}

	if err := s.lr.SaveToken(userId, token); err != nil {
		return "", "", err
	}

	return token, expiresIn, nil
}

func (s *loginService) UpdateToken(userId string) (token string, expiresIn string, err error) {
	token, expiresIn, err = createJWT(userId)
	if err != nil {
		return "", "", err
	}

	if err := s.lr.UpdateToken(userId, token); err != nil {
		return "", "", err
	}

	return token, expiresIn, nil
}

func createJWT(userId string) (token string, expiresIn string, err error) {
	const EXPIRES_IN = 24 // 24時間
	expireTime := time.Now().Add(time.Hour * EXPIRES_IN)
	secretKey := os.Getenv("SECRET_KEY")

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userId,
			"exp": expireTime.Unix(),
		})
	token, err = t.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return token, expireTime.Format(time.RFC3339), nil
}