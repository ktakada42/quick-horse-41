package login

import (
	"errors"
	"net/http"

	"app/user"
	"app/utils"
)

type useCaseInterface interface {
	login(email, password string) (token string, expiresIn string, err error)
}

type useCaseStruct struct {
	ls serviceInterface
	lr repositoryInterface
	ur user.RepositoryInterface
}

func NewLoginUseCase(ls serviceInterface, lr repositoryInterface, ur user.RepositoryInterface) useCaseInterface {
	return &useCaseStruct{ls: ls, lr: lr, ur: ur}
}

func (u *useCaseStruct) login(email, password string) (token string, expiresIn string, err error) {
	userId, err := u.ur.GetUserIdByEmail(email)
	if err != nil {
		return "", "", err
	}
	if userId == "" {
		// メールアドレスが間違っているのかパスワードが間違っているのかを知らせないために、Unauthorizedを返す
		return "", "", utils.NewHttpError(http.StatusUnauthorized, errors.New(utils.LoginErrorMessage))
	}

	isCorrect, err := u.ls.isPasswordCorrect(userId, password)
	if err != nil {
		return "", "", err
	}
	if !isCorrect {
		// メールアドレスが間違っているのかパスワードが間違っているのかを知らせないために、Unauthorizedを返す
		return "", "", utils.NewHttpError(http.StatusUnauthorized, errors.New(utils.LoginErrorMessage))
	}

	token, err = u.lr.getToken(userId)
	if err != nil {
		return "", "", err
	}
	if token == "" {
		return u.ls.createToken(userId)
	}

	return u.ls.updateToken(userId)
}
