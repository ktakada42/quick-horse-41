package login

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"app/user"
)

type serviceInterface interface {
	isPasswordCorrect(userId, password string) (bool, error)
	createToken(userId string) (token string, expiresIn string, err error)
	updateToken(userId string) (token string, expiresIn string, err error)
}

type serviceStruct struct {
	lr repositoryInterface
	ur user.RepositoryInterface
}

func NewLoginService(lr repositoryInterface, ur user.RepositoryInterface) serviceInterface {
	return &serviceStruct{lr: lr, ur: ur}
}

func (s *serviceStruct) isPasswordCorrect(userId, password string) (bool, error) {
	hashedPassword, err := s.ur.GetHashedPassword(userId)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func (s *serviceStruct) createToken(userId string) (token string, expiresIn string, err error) {
	token, expiresIn, err = createJWT(userId)
	if err != nil {
		return "", "", err
	}

	if err := s.lr.createToken(userId, token); err != nil {
		return "", "", err
	}

	return token, expiresIn, nil
}

func (s *serviceStruct) updateToken(userId string) (token string, expiresIn string, err error) {
	token, expiresIn, err = createJWT(userId)
	if err != nil {
		return "", "", err
	}

	if err := s.lr.updateToken(userId, token); err != nil {
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
