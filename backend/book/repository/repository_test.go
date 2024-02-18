package repository

import (
	"app/book/entity"
	"app/book/repository/testdata"
	"context"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBuildCanBorrowSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			want: `
        AND NOT EXISTS (
            SELECT 1 FROM borrowed_book bb
            WHERE
                ob.book_id = bb.book_id
                AND ob.office_id = bb.office_id
                AND ob.office_book_id = bb.office_book_id
        )
      `,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			r := &repositoryStruct{db: db}

			// 生の文字列リテラルを比較するため、Fieldsを使ってスペースを削除してから比較する
			query := strings.Join(strings.Fields(r.buildCanBorrowSQL()), " ")
			want := strings.Join(strings.Fields(tt.want), " ")
			if query != want {
				t.Errorf("buildCanBorrowSQL() query = %v, want %v", query, want)
			}
		})
	}
}

func TestBuildQSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			want: `
        AND (
            bm.title LIKE ?
            OR bm.author LIKE ?
            OR bm.publisher LIKE ?
        )
      `,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			r := &repositoryStruct{db: db}

			query := strings.Join(strings.Fields(r.buildQSQL()), " ")
			want := strings.Join(strings.Fields(tt.want), " ")
			if query != want {
				t.Errorf("buildQSQL() query = %v, want %v", query, want)
			}
		})
	}
}

func TestBuildSortSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		sort    string
		want    string
		wantErr bool
	}{
		{
			name:    "OK: reg-date-desc",
			sort:    "reg-date-desc",
			want:    "ORDER BY MAX(ob.reg_date) DESC ",
			wantErr: false,
		},
		{
			name:    "OK: reg-date-asc",
			sort:    "reg-date-asc",
			want:    "ORDER BY MAX(ob.reg_date) ASC ",
			wantErr: false,
		},
		{
			name:    "OK: rating-desc",
			sort:    "rating-desc",
			want:    "ORDER BY AVG(r.rating) DESC ",
			wantErr: false,
		},
		{
			name:    "OK: rating-asc",
			sort:    "rating-asc",
			want:    "ORDER BY AVG(r.rating) ASC ",
			wantErr: false,
		},
		{
			name:    "OK: review-count-desc",
			sort:    "review-count-desc",
			want:    "ORDER BY COUNT(r.review_id) DESC ",
			wantErr: false,
		},
		{
			name:    "OK: review-count-asc",
			sort:    "review-count-asc",
			want:    "ORDER BY COUNT(r.review_id) ASC ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			r := &repositoryStruct{db: db}

			if r.buildSortSQL(tt.sort) != tt.want {
				t.Errorf("buildQSQL() r.buildSortSQL(tt.sort) = %v, want %v", r.buildSortSQL(tt.sort), tt.want)
			}
		})
	}
}

func TestGetOfficeBooks(t *testing.T) {
	t.Parallel()

	baseSQL := `
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
	canBorrowSQL := `
		AND NOT EXISTS (
					SELECT 1 FROM borrowed_book bb
					WHERE
							ob.book_id = bb.book_id
							AND ob.office_id = bb.office_id
							AND ob.office_book_id = bb.office_book_id
		)
	`
	qSQL := `
		AND (
				bm.title LIKE ?
				OR bm.author LIKE ?
				OR bm.publisher LIKE ?
		)
	`
	officeIDSQL := "AND ob.office_id = ? "
	groupBySQL := "GROUP BY ob.book_id "
	limitOffsetSQL := "LIMIT ? OFFSET ?"

	tests := []struct {
		name        string
		in          GetOfficeBooksRepositoryInput
		mockClosure func(mock sqlmock.Sqlmock)
		want        []entity.BookForList
		wantErr     bool
	}{
		{
			name: "OK",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       20,
				OfficeID:    "ofN1A123456789A123456789A1",
				Q:           "ハリー",
				CanBorrow:   true,
				IsAvailable: true,
				Sort:        "reg-date-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+canBorrowSQL+qSQL+officeIDSQL+groupBySQL+"ORDER BY MAX(ob.reg_date) DESC "+limitOffsetSQL,
					),
				).
					WithArgs(true, "%ハリー%", "%ハリー%", "%ハリー%", "ofN1A123456789A123456789A1", 20, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[1].BookID,
								testdata.BookForListTestData[1].Title,
								testdata.BookForListTestData[1].CoverID,
								testdata.BookForListTestData[1].LastRegDate,
								testdata.BookForListTestData[1].ReviewCount,
								testdata.BookForListTestData[1].Rating,
							).
							AddRow(
								testdata.BookForListTestData[0].BookID,
								testdata.BookForListTestData[0].Title,
								testdata.BookForListTestData[0].CoverID,
								testdata.BookForListTestData[0].LastRegDate,
								testdata.BookForListTestData[0].ReviewCount,
								testdata.BookForListTestData[0].Rating,
							),
					)
			},
			want: []entity.BookForList{
				testdata.BookForListTestData[1],
				testdata.BookForListTestData[0],
			},
			wantErr: false,
		},
		{
			name: "OK: 貸出可能な本以外も取得",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       1,
				OfficeID:    "ofN1A123456789A123456789A1",
				Q:           "ハリー",
				CanBorrow:   false,
				IsAvailable: true,
				Sort:        "reg-date-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+qSQL+officeIDSQL+groupBySQL+"ORDER BY MAX(ob.reg_date) DESC "+limitOffsetSQL,
					),
				).
					WithArgs(true, "%ハリー%", "%ハリー%", "%ハリー%", "ofN1A123456789A123456789A1", 1, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[1].BookID,
								testdata.BookForListTestData[1].Title,
								testdata.BookForListTestData[1].CoverID,
								testdata.BookForListTestData[1].LastRegDate,
								testdata.BookForListTestData[1].ReviewCount,
								testdata.BookForListTestData[1].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[1]},
			wantErr: false,
		},
		{
			name: "OK: office ID, Q のキーワードなし",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       20,
				OfficeID:    "",
				Q:           "",
				CanBorrow:   false,
				IsAvailable: true,
				Sort:        "reg-date-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+groupBySQL+"ORDER BY MAX(ob.reg_date) DESC "+limitOffsetSQL,
					),
				).
					WithArgs(true, 20, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[2].BookID,
								testdata.BookForListTestData[2].Title,
								testdata.BookForListTestData[2].CoverID,
								testdata.BookForListTestData[2].LastRegDate,
								testdata.BookForListTestData[2].ReviewCount,
								testdata.BookForListTestData[2].Rating,
							).
							AddRow(
								testdata.BookForListTestData[1].BookID,
								testdata.BookForListTestData[1].Title,
								testdata.BookForListTestData[1].CoverID,
								testdata.BookForListTestData[1].LastRegDate,
								testdata.BookForListTestData[1].ReviewCount,
								testdata.BookForListTestData[1].Rating,
							).
							AddRow(
								testdata.BookForListTestData[0].BookID,
								testdata.BookForListTestData[0].Title,
								testdata.BookForListTestData[0].CoverID,
								testdata.BookForListTestData[0].LastRegDate,
								testdata.BookForListTestData[0].ReviewCount,
								testdata.BookForListTestData[0].Rating,
							),
					)
			},
			want: []entity.BookForList{
				testdata.BookForListTestData[2],
				testdata.BookForListTestData[1],
				testdata.BookForListTestData[0],
			},
			wantErr: false,
		},
		{
			name: "OK: sort reg-date-asc",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       1,
				OfficeID:    "",
				Q:           "",
				CanBorrow:   false,
				IsAvailable: true,
				Sort:        "reg-date-asc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+groupBySQL+"ORDER BY MAX(ob.reg_date) ASC "+limitOffsetSQL,
					),
				).
					WithArgs(true, 1, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[0].BookID,
								testdata.BookForListTestData[0].Title,
								testdata.BookForListTestData[0].CoverID,
								testdata.BookForListTestData[0].LastRegDate,
								testdata.BookForListTestData[0].ReviewCount,
								testdata.BookForListTestData[0].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[0]},
			wantErr: false,
		},
		{
			name: "OK: sort rating-desc",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       1,
				OfficeID:    "",
				Q:           "",
				CanBorrow:   false,
				IsAvailable: true,
				Sort:        "rating-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+groupBySQL+"ORDER BY AVG(r.rating) DESC "+limitOffsetSQL,
					),
				).
					WithArgs(true, 1, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[1].BookID,
								testdata.BookForListTestData[1].Title,
								testdata.BookForListTestData[1].CoverID,
								testdata.BookForListTestData[1].LastRegDate,
								testdata.BookForListTestData[1].ReviewCount,
								testdata.BookForListTestData[1].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[1]},
			wantErr: false,
		},
		{
			name: "OK: sort rating-asc",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       1,
				OfficeID:    "",
				Q:           "",
				CanBorrow:   false,
				IsAvailable: true,
				Sort:        "rating-asc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+groupBySQL+"ORDER BY AVG(r.rating) ASC "+limitOffsetSQL,
					),
				).
					WithArgs(true, 1, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[2].BookID,
								testdata.BookForListTestData[2].Title,
								testdata.BookForListTestData[2].CoverID,
								testdata.BookForListTestData[2].LastRegDate,
								testdata.BookForListTestData[2].ReviewCount,
								testdata.BookForListTestData[2].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[2]},
			wantErr: false,
		},
		{
			name: "OK: sort review-count-desc",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       1,
				OfficeID:    "",
				Q:           "",
				CanBorrow:   false,
				IsAvailable: true,
				Sort:        "review-count-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+groupBySQL+"ORDER BY COUNT(r.review_id) DESC "+limitOffsetSQL,
					),
				).
					WithArgs(true, 1, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[1].BookID,
								testdata.BookForListTestData[1].Title,
								testdata.BookForListTestData[1].CoverID,
								testdata.BookForListTestData[1].LastRegDate,
								testdata.BookForListTestData[1].ReviewCount,
								testdata.BookForListTestData[1].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[1]},
			wantErr: false,
		},
		{
			name: "OK: sort review-count-asc",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       1,
				OfficeID:    "",
				Q:           "",
				CanBorrow:   false,
				IsAvailable: true,
				Sort:        "review-count-asc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						baseSQL+groupBySQL+"ORDER BY COUNT(r.review_id) ASC "+limitOffsetSQL),
				).
					WithArgs(true, 1, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[2].BookID,
								testdata.BookForListTestData[2].Title,
								testdata.BookForListTestData[2].CoverID,
								testdata.BookForListTestData[2].LastRegDate,
								testdata.BookForListTestData[2].ReviewCount,
								testdata.BookForListTestData[2].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[2]},
			wantErr: false,
		},
		{
			name: "OK: limit",
			in: GetOfficeBooksRepositoryInput{
				Offset:      0,
				Limit:       1,
				OfficeID:    "ofN1A123456789A123456789A1",
				Q:           "ハリー",
				CanBorrow:   true,
				IsAvailable: true,
				Sort:        "reg-date-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					baseSQL+canBorrowSQL+qSQL+officeIDSQL+groupBySQL+"ORDER BY MAX(ob.reg_date) DESC "+limitOffsetSQL),
				).
					WithArgs(true, "%ハリー%", "%ハリー%", "%ハリー%", "ofN1A123456789A123456789A1", 1, 0).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[1].BookID,
								testdata.BookForListTestData[1].Title,
								testdata.BookForListTestData[1].CoverID,
								testdata.BookForListTestData[1].LastRegDate,
								testdata.BookForListTestData[1].ReviewCount,
								testdata.BookForListTestData[1].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[1]},
			wantErr: false,
		},
		{
			name: "OK: offset",
			in: GetOfficeBooksRepositoryInput{
				Offset:      1,
				Limit:       1,
				OfficeID:    "ofN1A123456789A123456789A1",
				Q:           "ハリー",
				CanBorrow:   true,
				IsAvailable: true,
				Sort:        "reg-date-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					baseSQL+canBorrowSQL+qSQL+officeIDSQL+groupBySQL+"ORDER BY MAX(ob.reg_date) DESC "+limitOffsetSQL),
				).
					WithArgs(true, "%ハリー%", "%ハリー%", "%ハリー%", "ofN1A123456789A123456789A1", 1, 1).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}).
							AddRow(
								testdata.BookForListTestData[0].BookID,
								testdata.BookForListTestData[0].Title,
								testdata.BookForListTestData[0].CoverID,
								testdata.BookForListTestData[0].LastRegDate,
								testdata.BookForListTestData[0].ReviewCount,
								testdata.BookForListTestData[0].Rating,
							),
					)
			},
			want:    []entity.BookForList{testdata.BookForListTestData[0]},
			wantErr: false,
		},
		{
			name: "OK: offsetが大きいため該当の本なし",
			in: GetOfficeBooksRepositoryInput{
				Offset:      20,
				Limit:       1,
				OfficeID:    "ofN1A123456789A123456789A1",
				Q:           "ハリー",
				CanBorrow:   true,
				IsAvailable: true,
				Sort:        "reg-date-desc",
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					baseSQL+canBorrowSQL+qSQL+officeIDSQL+groupBySQL+"ORDER BY MAX(ob.reg_date) DESC "+limitOffsetSQL),
				).
					WithArgs(true, "%ハリー%", "%ハリー%", "%ハリー%", "ofN1A123456789A123456789A1", 1, 20).
					WillReturnRows(
						sqlmock.NewRows([]string{"bookID", "title", "coverId", "lastRegDate", "reviewCount", "rating"}),
					)
			},
			want:    []entity.BookForList{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockClosure(mock)

			r := NewBookRepository(db)

			got, err := r.GetOfficeBooks(ctx, tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOfficeBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOfficeBooks() got = %v, want %v", got, tt.want)
			}
		})
	}
}
