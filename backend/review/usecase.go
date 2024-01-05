package review

type useCaseInterface interface {
	reviews(offset int, limit int) ([]Review, error)
}

type useCaseStruct struct {
	rr repositoryInterface
}

func NewReviewUseCase(rr repositoryInterface) useCaseInterface {
	return &useCaseStruct{rr: rr}
}

func (u *useCaseStruct) reviews(offset int, limit int) ([]Review, error) {
	reviews, err := u.rr.GetReviews(offset, limit)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
