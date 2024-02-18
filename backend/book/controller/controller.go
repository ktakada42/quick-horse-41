package controller

import (
	"app/book/usecase"
	"app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/form"
)

var validSorts = map[string]bool{
	"reg-date-desc":     true,
	"reg-date-asc":      true,
	"rating-desc":       true,
	"rating-asc":        true,
	"review-count-desc": true,
	"review-count-asc":  true,
}

type getOfficeBooksQueryParameters struct {
	Offset      uint32 `form:"offset"`
	Limit       uint8  `form:"limit"`
	OfficeID    string `form:"office-id"`
	Q           string `form:"q"`
	CanBorrow   bool   `form:"can-borrow"`
	IsAvailable bool   `form:"is-available"`
	Sort        string `form:"sort"`
}

type ControllerInterface interface {
	GetOfficeBooksController(w http.ResponseWriter, r *http.Request)
}

type controllerStruct struct {
	bu usecase.UseCaseInterface
}

func NewBookController(bu usecase.UseCaseInterface) ControllerInterface {
	return &controllerStruct{bu: bu}
}

func (c *controllerStruct) getOfficeBooksValidateRequest(r *http.Request) (getOfficeBooksQueryParameters, error) {
	var queryParams getOfficeBooksQueryParameters
	var errorMessages []string

	if r.Method != http.MethodGet {
		return queryParams, fmt.Errorf("method: %s is invalid", r.Method)
	}

	decoder := form.NewDecoder()
	err := decoder.Decode(&queryParams, r.URL.Query())
	if err != nil {
		return queryParams, fmt.Errorf("decode error: %v", err)
	}

	// check required parameters
	if _, ok := r.URL.Query()["offset"]; !ok {
		errorMessages = append(errorMessages, "offset")
	}
	if _, ok := r.URL.Query()["limit"]; !ok {
		errorMessages = append(errorMessages, "limit")
	}
	if _, ok := r.URL.Query()["can-borrow"]; !ok {
		errorMessages = append(errorMessages, "can-borrow")
	}
	if _, ok := r.URL.Query()["is-available"]; !ok {
		errorMessages = append(errorMessages, "is-available")
	}
	if _, ok := r.URL.Query()["sort"]; !ok {
		errorMessages = append(errorMessages, "sort")
	}

	if len(errorMessages) > 0 {
		message := fmt.Sprintf("Missing required query parameters: %s.", strings.Join(errorMessages, ", "))
		return queryParams, fmt.Errorf(message)
	}

	// validate parameters
	if queryParams.Offset < 0 {
		errorMessages = append(errorMessages, "offset")
	}
	if queryParams.Limit < 0 || queryParams.Limit > 100 {
		errorMessages = append(errorMessages, "limit")
	}
	if queryParams.OfficeID != "" && len(queryParams.OfficeID) != 26 {
		errorMessages = append(errorMessages, "office-id")
	}
	if utf8.RuneCountInString(queryParams.Q) > 100 {
		errorMessages = append(errorMessages, "q")
	}
	if _, exists := validSorts[queryParams.Sort]; !exists {
		errorMessages = append(errorMessages, "sort")
	}

	if len(errorMessages) > 0 {
		message := fmt.Sprintf("Invalid query parameters: %s.", strings.Join(errorMessages, ", "))
		return queryParams, fmt.Errorf(message)
	}

	return queryParams, nil
}

// 本一覧取得
func (c *controllerStruct) GetOfficeBooksController(w http.ResponseWriter, r *http.Request) {
	queryParams, err := c.getOfficeBooksValidateRequest(r)

	if err != nil {
		utils.SetJsonError(w, err, http.StatusBadRequest)
		return
	}

	response, err := c.bu.GetOfficeBooks(r.Context(), usecase.GetOfficeBooksUsecaseInput{
		Offset:      queryParams.Offset,
		Limit:       queryParams.Limit,
		OfficeID:    queryParams.OfficeID,
		Q:           queryParams.Q,
		CanBorrow:   queryParams.CanBorrow,
		IsAvailable: queryParams.IsAvailable,
		Sort:        queryParams.Sort,
	})
	if err != nil {
		var httpError *utils.HttpError
		if errors.As(err, &httpError) {
			utils.SetJsonError(w, err, httpError.GetStatusCode())
		} else {
			utils.SetJsonError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
