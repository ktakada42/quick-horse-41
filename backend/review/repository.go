package review

import (
	"database/sql"
)

type repositoryInterface interface {
	GetReviews(offset int, limit int) ([]Review, error)
}

type repositoryStruct struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) repositoryInterface {
	return &repositoryStruct{db: db}
}

func (r *repositoryStruct) GetReviews(offset int, limit int) ([]Review, error) {
	rows, err := r.db.Query("SELECT * from review LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	reviews := []Review{}
	for rows.Next() {
		var review Review
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
