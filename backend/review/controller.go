package review

import (
	"app/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ControllerInterface interface {
	GetReviews(w http.ResponseWriter, r *http.Request)
}

type controllerStruct struct {
	ruc useCaseInterface
}

func NewReviewController(ruc useCaseInterface) ControllerInterface {
	return &controllerStruct{ruc: ruc}
}

type Request struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type review struct {
	BookID   string    `json:"bookId"`
	ReviewID int       `json:"reviewId"`
	UserID   string    `json:"userId"`
	Rating   int8      `json:"rating"`
	Review   string    `json:"review"`
	RegDate  time.Time `json:"regDate"`
}

type Response struct {
	Reviews []review `json:"reviews"`
}

// 全てのレビューを取得する
func (c *controllerStruct) GetReviews(w http.ResponseWriter, r *http.Request) {
	// リクエストボディを取得
	var body Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		if err == io.EOF {
			utils.SetJsonError(w, fmt.Errorf("request body is empty"), http.StatusBadRequest)
		} else {
			utils.SetJsonError(w, fmt.Errorf("request body is invalid: %w", err), http.StatusBadRequest)
		}
		return
	}
	offset, limit := &body.Offset, &body.Limit

	// リクエストボディのバリデーション
	var errMessages []string
	if *offset < 0 {
		errMessages = append(errMessages, "offset must be non-negative")
	}
	if *limit < 1 || *limit > 100 {
		errMessages = append(errMessages, "limit must be between 1 and 100")
	}
	if len(errMessages) > 0 {
		utils.SetJsonError(w, fmt.Errorf(strings.Join(errMessages, ", ")), http.StatusBadRequest)
		return
	}

	// usecaseの呼び出し
	reviews, err := c.ruc.GetReviews(*offset, *limit)
	if err != nil {
		utils.SetJsonError(w, err, http.StatusInternalServerError)
		return
	}

	// responseに変換
	var response Response
	response.Reviews = []review{}
	for _, item := range reviews {
		response.Reviews = append(response.Reviews, review(item))
	}

	// レスポンスの返却
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
