package review

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func setTestData() ([]Review, *sqlmock.Rows) {
	reviews := []Review{
		{"book1", 1, "user1", 5, "素晴らしい本でした。", timeParser("2023-01-01 00:00:00")},
		{"book2", 2, "user2", 4, "面白かったです。", timeParser("2023-01-02 00:00:00")},
	}

	rows := sqlmock.NewRows([]string{"book_id", "review_id", "user_id", "rating", "review", "reg_date"}).
		AddRow("book1", 1, "user1", 5, "素晴らしい本でした。", timeParser("2023-01-01 00:00:00")).
		AddRow("book2", 2, "user2", 4, "面白かったです。", timeParser("2023-01-02 00:00:00"))

	return reviews, rows
}

func timeParser(timeStr string) time.Time {
	layout := "2006-01-02 15:04:05"
	time, _ := time.Parse(layout, timeStr)
	return time
}
func TestRepositoryGetReviews(t *testing.T) {
	t.Parallel()

	// テストデータの準備
	reviews, rows := setTestData()
	query := "SELECT \\* from review LIMIT \\? OFFSET \\?"

	tests := map[string]struct {
		offset      int
		limit       int
		mockClosure func(mock sqlmock.Sqlmock)
		want        []Review
		wantErr     bool
	}{
		"正常系: offset = 0, limit = 10": {
			offset: 0,
			limit:  10,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			want:    reviews,
			wantErr: false,
		},
		"異常系: offset が負の値": {
			offset: -1,
			limit:  10,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(10, -1).
					WillReturnError(fmt.Errorf("error"))
			},
			want:    nil,
			wantErr: true,
		},
		"異常系: limit が 1 より小さい": {
			offset: 0,
			limit:  -1,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(-1, 0).
					WillReturnError(fmt.Errorf("error"))
			},
			want:    nil,
			wantErr: true,
		},
		"異常系: limit が 100 より大きい": {
			offset: 0,
			limit:  1000,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(1000, 0).
					WillReturnError(fmt.Errorf("error"))
			},
			want:    nil,
			wantErr: true,
		},
		"異常系: offset と limit が範囲外": {
			offset: -1,
			limit:  -1,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(-1, -1).
					WillReturnError(fmt.Errorf("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		// ttとnameの値をコピーして使用する
		tt := tt
		name := name

		// テストの実行
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Logf("test name: %s, start: %s, address of tt: %p", name, time.Now(), &tt)

			// sqlmockの初期化
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockClosure(mock)

			// repositoryの初期化
			repository := NewReviewRepository(db)

			// テスト対象の実行
			got, err := repository.GetReviews(tt.offset, tt.limit)
			t.Logf("test name: %s, got: %v", name, got)

			// 結果の検証
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("error = %v, want = %v", got, tt.want)
			}
		})
	}
}
