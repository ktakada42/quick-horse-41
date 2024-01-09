package usecase

import (
	"errors"
	"net/http"

	lr "app/login/repository"
	"app/login/service"
	ur "app/user/repository"
	"app/utils"
)

type LoginUseCase interface {
	Login(email, password string) (token string, expiresIn string, err error)
}

type loginUseCase struct {
	ls service.LoginService
	lr lr.LoginRepository
	ur ur.UserRepository
}

func NewLoginUseCase(ls service.LoginService, lr lr.LoginRepository, ur ur.UserRepository) LoginUseCase {
	return &loginUseCase{ls: ls, lr: lr, ur: ur}
}

func (u *loginUseCase) Login(email, password string) (token string, expiresIn string, err error) {
	userId, err := u.ur.GetUserIdByEmail(email)
	if err != nil {
		return "", "", err
	}
	if userId == "" {
		// メールアドレスが間違っているのかパスワードが間違っているのかを知らせないために、Unauthorizedを返す
		return "", "", utils.NewHttpError(http.StatusUnauthorized, errors.New(utils.LoginErrorMessage))
	}

	isCorrect, err := u.ls.IsPasswordCorrect(userId, password)
	if err != nil {
		return "", "", err
	}
	if !isCorrect {
		// メールアドレスが間違っているのかパスワードが間違っているのかを知らせないために、Unauthorizedを返す
		return "", "", utils.NewHttpError(http.StatusUnauthorized, errors.New(utils.LoginErrorMessage))
	}

	token, err = u.lr.GetToken(userId)
	if err != nil {
		return "", "", err
	}
	if token == "" {
		return u.ls.CreateToken(userId)
	}

	return u.ls.UpdateToken(userId)
}
