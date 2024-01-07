package review

import (
	"errors"
	"log"
)

type useCaseInterface interface {
	GetReviews(offset int, limit int) ([]Review, error)
}

type useCaseStruct struct {
	rr repositoryInterface
}

func NewReviewUseCase(rr repositoryInterface) useCaseInterface {
	return &useCaseStruct{rr: rr}
}

// 全てのレビューを取得する
func (uc *useCaseStruct) GetReviews(offset int, limit int) ([]Review, error) {
	reviews, err := uc.rr.GetReviews(offset, limit)
	if err != nil {
		// エラーの詳細はログに出力する
		log.Printf("review usecase GetReviews() error: %v", err)
		return nil, errors.New("internal server error")
	}
	return reviews, nil
}
