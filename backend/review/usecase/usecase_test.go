package usecase

import (
	"app/review/entity"
	"app/review/mock/mock_review"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func setTestData() []entity.Review {
	reviews := []entity.Review{
		{
			BookID:   "book1",
			ReviewID: 1,
			UserID:   "user1",
			Rating:   5,
			Review:   "素晴らしい本でした。",
			RegDate:  timeParser("2023-01-01 00:00:00"),
		},
		{
			BookID:   "book2",
			ReviewID: 2,
			UserID:   "user2",
			Rating:   4,
			Review:   "面白かったです。",
			RegDate:  timeParser("2023-01-02 00:00:00"),
		},
	}
	return reviews
}

func timeParser(timeStr string) time.Time {
	layout := "2006-01-02 15:04:05"
	time, _ := time.Parse(layout, timeStr)
	return time
}

type mocks struct {
	reviewRepository *mock_review.MockRepositoryInterface
}

func newWithMocks(t *testing.T) (UseCaseInterface, *mocks) {
	ctrl := gomock.NewController(t)
	mockReviewRepository := mock_review.NewMockRepositoryInterface(ctrl)
	ruc := NewReviewUseCase(mockReviewRepository)
	return ruc, &mocks{
		reviewRepository: mockReviewRepository,
	}
}

func TestGetReviews(t *testing.T) {
	t.Parallel()

	// テストデータの準備
	reviews := setTestData()

	type args struct {
		offset int
		limit  int
	}

	type want struct {
		res []entity.Review
		err error
	}

	tests := map[string]struct {
		args    args
		want    want
		wantErr bool
	}{
		"正常系: offset = 0, limit = 10": {
			args: args{
				offset: 0,
				limit:  10,
			},
			want: want{
				res: reviews,
				err: nil,
			},
			wantErr: false,
		},
		"異常系: offset が負の値": {
			args: args{
				offset: -1,
				limit:  10,
			},
			want: want{
				res: nil,
				err: errors.New("error"),
			},
			wantErr: true,
		},
		"異常系: limit が 1 より小さい": {
			args: args{
				offset: 0,
				limit:  -1,
			},
			want: want{
				res: nil,
				err: errors.New("error"),
			},
			wantErr: true,
		},
		"異常系: limit が 100 より大きい": {
			args: args{
				offset: 0,
				limit:  1000,
			},
			want: want{
				res: nil,
				err: errors.New("error"),
			},
			wantErr: true,
		},
		"異常系: offset と limit が範囲外": {
			args: args{
				offset: -1,
				limit:  -1,
			},
			want: want{
				res: nil,
				err: errors.New("error"),
			},
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

			// 初期化
			ruc, m := newWithMocks(t)

			// モックの設定
			m.reviewRepository.EXPECT().
				GetReviews(tt.args.offset, tt.args.limit).
				Return(tt.want.res, tt.want.err)

			// テスト対象の実行
			got, err := ruc.GetReviews(tt.args.offset, tt.args.limit)
			t.Logf("test name: %s, got: %v", name, got)

			// 結果の検証
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want.res) {
				t.Errorf("error = %v, want = %v", got, tt.want)
			}
		})
	}
}
