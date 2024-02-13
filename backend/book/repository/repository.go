package repository

import (
	"app/book/entity"
	"database/sql"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type RepositoryInterface interface {
	GetOfficeBooks(offset uint32, limit uint8, officeId string, q string, canBorrow bool, isAvailable bool, sort string) ([]entity.BookForList, error)
}

type repositoryStruct struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) RepositoryInterface {
	return &repositoryStruct{db: db}
}

func (r *repositoryStruct) GetOfficeBooks(offset uint32, limit uint8, officeId string, q string, canBorrow bool, isAvailable bool, sort string) ([]entity.BookForList, error) {
	var args []interface{}

	strQuery := `
    SELECT
        ob.book_id as bookID
      , bm.title as title
      , bm.cover_id as coverId
      , max(ob.reg_date) as lastRegDate
      , count(r.review_id) as reviewCount
      , avg(r.rating) as rating
    FROM
      	office_book ob
    LEFT JOIN
      	book_master bm ON ob.book_id = bm.book_id
    LEFT JOIN
      	review r ON ob.book_id = r.book_id
		WHERE
				ob.available = ?
  `
	args = append(args, isAvailable)

	// 貸出可能な本のみ取得するか
	if canBorrow {
		strQuery += `
			AND NOT EXISTS (
        	SELECT 1 FROM borrowed_book bb
        	WHERE
            	ob.book_id = bb.book_id
            	AND ob.office_id = bb.office_id
            	AND ob.office_book_id = bb.office_book_id
    	)
		`
	}

	// キーワード検索
	if q != "" {
		q = "%" + q + "%"
		strQuery += `
			AND (
					bm.title like ?
					or bm.author like ?
					or bm.publisher like ?
			)
		`
		args = append(args, q, q, q)
	}

	// オフィスIDの指定
	if officeId != "" {
		strQuery += "AND ob.office_id = ? "
		args = append(args, officeId)
	}

	// グループ化
	strQuery += "GROUP BY ob.book_id "

	// ソート
	strQuery += "ORDER BY "
	if sort == "reg-date-desc" {
		strQuery += "MAX(ob.reg_date) DESC "
	} else if sort == "reg-date-asc" {
		strQuery += "MAX(ob.reg_date) ASC "
	} else if sort == "rating-desc" {
		strQuery += "AVG(r.rating) DESC "
	} else if sort == "rating-asc" {
		strQuery += "AVG(r.rating) ASC "
	} else if sort == "review-count-desc" {
		strQuery += "COUNT(r.review_id) DESC "
	} else if sort == "review-count-asc" {
		strQuery += "COUNT(r.review_id) ASC "
	}

	// limit
	strQuery += "LIMIT ? "
	args = append(args, limit)

	// offset
	strQuery += "OFFSET ? "
	args = append(args, offset)

	rows, err := r.db.Query(strQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookForList := []entity.BookForList{}
	for rows.Next() {
		var book entity.BookForList
		if err := rows.Scan(&book.BookId, &book.Title, &book.CoverId, &book.LastRegDate, &book.ReviewCount, &book.Rating); err != nil {
			return nil, err
		}
		bookForList = append(bookForList, book)
	}

	return bookForList, nil
}
