package repository

import (
	"app/review/entity"
	"database/sql"
)

type RepositoryInterface interface {
	GetReviews(offset int, limit int) ([]entity.Review, error)
}

type repositoryStruct struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) RepositoryInterface {
	return &repositoryStruct{db: db}
}

// 全てのレビューを取得する
func (r *repositoryStruct) GetReviews(offset int, limit int) ([]entity.Review, error) {
	rows, err := r.db.Query("SELECT * from review LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	return convertToReviews(rows)
}

// rows型をReview型に紐づける
func convertToReviews(rows *sql.Rows) ([]entity.Review, error) {
	reviews := []entity.Review{}
	for rows.Next() {
		var review entity.Review
		if err := rows.Scan(&review.BookID, &review.ReviewID, &review.UserID, &review.Rating, &review.Review, &review.RegDate); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}
