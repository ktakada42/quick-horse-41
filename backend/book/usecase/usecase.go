package usecase

import (
	"app/book/entity"
	"app/book/repository"

	"log"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type UseCaseInterface interface {
	GetOfficeBooks(offset uint32, limit uint8, officeId string, q string, canBorrow bool, isAvailable bool, sort string) ([]entity.BookForList, error)
}

type useCaseStruct struct {
	br repository.RepositoryInterface
}

func NewBookUseCase(br repository.RepositoryInterface) UseCaseInterface {
	return &useCaseStruct{br: br}
}

func (u *useCaseStruct) GetOfficeBooks(offset uint32, limit uint8, officeId string, q string, canBorrow bool, isAvailable bool, sort string) ([]entity.BookForList, error) {
	books, err := u.br.GetOfficeBooks(offset, limit, officeId, q, canBorrow, isAvailable, sort)
	if err != nil {
		log.Printf("book usercase GetOfficeBooks() error: %v", err)
		return nil, err
	}

	return books, nil
}
