package usecase

import (
	"app/book/entity"
	"app/book/repository"
	"context"

	"log"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type UseCaseInterface interface {
	GetOfficeBooks(ctx context.Context, in GetOfficeBooksUsecaseInput) ([]entity.BookForList, error)
}

type useCaseStruct struct {
	br repository.RepositoryInterface
}

type GetOfficeBooksUsecaseInput struct {
	Offset      uint32
	Limit       uint8
	OfficeID    string
	Q           string
	CanBorrow   bool
	IsAvailable bool
	Sort        string
}

func NewBookUseCase(br repository.RepositoryInterface) UseCaseInterface {
	return &useCaseStruct{br: br}
}

func (u *useCaseStruct) GetOfficeBooks(ctx context.Context, in GetOfficeBooksUsecaseInput) ([]entity.BookForList, error) {
	books, err := u.br.GetOfficeBooks(
		ctx,
		repository.GetOfficeBooksRepositoryInput{
			Offset:      in.Offset,
			Limit:       in.Limit,
			OfficeID:    in.OfficeID,
			Q:           in.Q,
			CanBorrow:   in.CanBorrow,
			IsAvailable: in.IsAvailable,
			Sort:        in.Sort,
		},
	)
	if err != nil {
		log.Printf("book usercase GetOfficeBooks() error: %v", err)
		return nil, err
	}

	return books, nil
}
