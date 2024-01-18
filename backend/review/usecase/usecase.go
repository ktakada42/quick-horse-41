package usecase

import (
	"app/review/entity"
	"app/review/repository"
	"errors"
	"log"
)

type UseCaseInterface interface {
	GetReviews(offset int, limit int) ([]entity.Review, error)
}

type useCaseStruct struct {
	rr repository.RepositoryInterface
}

func NewReviewUseCase(rr repository.RepositoryInterface) UseCaseInterface {
	return &useCaseStruct{rr: rr}
}

// 全てのレビューを取得する
func (uc *useCaseStruct) GetReviews(offset int, limit int) ([]entity.Review, error) {
	reviews, err := uc.rr.GetReviews(offset, limit)
	if err != nil {
		// エラーの詳細はログに出力する
		log.Printf("review usecase GetReviews() error: %v", err)
		return nil, errors.New("internal server error")
	}
	return reviews, nil
}
