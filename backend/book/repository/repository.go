package repository

import (
	"app/book/entity"
	"context"
	"database/sql"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type RepositoryInterface interface {
	GetOfficeBooks(ctx context.Context, in GetOfficeBooksRepositoryInput) ([]entity.BookForList, error)
}

type repositoryStruct struct {
	db *sql.DB
}

type GetOfficeBooksRepositoryInput struct {
	Offset      uint32
	Limit       uint8
	OfficeID    string
	Q           string
	CanBorrow   bool
	IsAvailable bool
	Sort        string
}

func NewBookRepository(db *sql.DB) RepositoryInterface {
	return &repositoryStruct{db: db}
}

func (r *repositoryStruct) buildCanBorrowSQL() string {
	return `
		AND NOT EXISTS (
				SELECT 1 FROM borrowed_book bb
				WHERE
						ob.book_id = bb.book_id
						AND ob.office_id = bb.office_id
						AND ob.office_book_id = bb.office_book_id
		)
	`
}

func (r *repositoryStruct) buildQSQL() string {
	return `
		AND (
				bm.title LIKE ?
				OR bm.author LIKE ?
				OR bm.publisher LIKE ?
		)
	`
}

func (r *repositoryStruct) buildSortSQL(sort string) string {
	query := "ORDER BY "

	if sort == "reg-date-desc" {
		query += "MAX(ob.reg_date) DESC "

	} else if sort == "reg-date-asc" {
		query += "MAX(ob.reg_date) ASC "

	} else if sort == "rating-desc" {
		query += "AVG(r.rating) DESC "

	} else if sort == "rating-asc" {
		query += "AVG(r.rating) ASC "

	} else if sort == "review-count-desc" {
		query += "COUNT(r.review_id) DESC "

	} else if sort == "review-count-asc" {
		query += "COUNT(r.review_id) ASC "
	}

	return query
}

func (r *repositoryStruct) GetOfficeBooks(ctx context.Context, in GetOfficeBooksRepositoryInput) ([]entity.BookForList, error) {
	var args []interface{}

	strQuery := `
    SELECT
        ob.book_id AS bookID
      , bm.title AS title
      , bm.cover_id AS coverId
      , MAX(ob.reg_date) AS lastRegDate
      , COUNT(r.review_id) AS reviewCount
      , AVG(r.rating) AS rating
    FROM
      	office_book ob
    LEFT JOIN
      	book_master bm ON ob.book_id = bm.book_id
    LEFT JOIN
      	review r ON ob.book_id = r.book_id
		WHERE
				ob.available = ?
  `
	args = append(args, in.IsAvailable)

	// 貸出可能な本のみ取得するか
	if in.CanBorrow {
		strQuery += r.buildCanBorrowSQL()
	}

	// キーワード検索
	if in.Q != "" {
		strQuery += r.buildQSQL()
		q := "%" + in.Q + "%"
		args = append(args, q, q, q)
	}

	// オフィスIDの指定
	if in.OfficeID != "" {
		// クエリを組み立てていくので、末尾にスペースを入れている
		strQuery += "AND ob.office_id = ? "
		args = append(args, in.OfficeID)
	}

	// グループ化
	strQuery += "GROUP BY ob.book_id "

	// ソート
	strQuery += r.buildSortSQL(in.Sort)

	// limit
	strQuery += "LIMIT ? "
	args = append(args, in.Limit)

	// offset
	strQuery += "OFFSET ? "
	args = append(args, in.Offset)

	rows, err := r.db.QueryContext(ctx, strQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookForList := []entity.BookForList{}
	for rows.Next() {
		var book entity.BookForList
		if err := rows.Scan(&book.BookID, &book.Title, &book.CoverID, &book.LastRegDate, &book.ReviewCount, &book.Rating); err != nil {
			return nil, err
		}
		bookForList = append(bookForList, book)
	}

	return bookForList, nil
}
