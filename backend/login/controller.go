package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"app/utils"
)

type controllerInterface interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type controllerStruct struct {
	lu useCaseInterface
}

func NewLoginController(lu useCaseInterface) controllerInterface {
	return &controllerStruct{lu: lu}
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Token     string `json:"token"`
	ExpiresIn string `json:"expiresIn"`
}

// nil以外を返すときはBad Requestを返すので、HttpErrorにはラップしない
func (c *controllerStruct) validateRequest(r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method: %s is invalid", r.Method)
	}

	if strings.ToLower(r.Header.Get("Content-Type")) != "application/json" {
		return fmt.Errorf("Content-Type: %s is invalid", r.Header.Get("Content-Type"))
	}

	return nil
}

func (c *controllerStruct) Login(w http.ResponseWriter, r *http.Request) {
	if err := c.validateRequest(r); err != nil {
		utils.SetJsonError(w, err, http.StatusBadRequest)
		return
	}

	var body Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.SetJsonError(w, fmt.Errorf("request body is invalid: %w", err), http.StatusBadRequest)
		return
	}
	email, password := body.Email, body.Password

	token, expiresIn, err := c.lu.login(email, password)
	if err != nil {
		var httpError *utils.HttpError
		if errors.As(err, &httpError) {
			utils.SetJsonError(w, err, httpError.GetStatusCode())
		} else {
			utils.SetJsonError(w, err, http.StatusInternalServerError)
		}
		return
	}

	response := Response{Token: token, ExpiresIn: expiresIn}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
