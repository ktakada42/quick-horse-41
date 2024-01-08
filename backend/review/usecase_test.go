package review

import (
	"app/mock/mock_review"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type mocks struct {
	reviewRepo *mock_review.MockrepositoryInterface
}

func newWithMocks(t *testing.T) (useCaseInterface, *mocks) {
	ctrl := gomock.NewController(t)
	mockReviewRepo := mock_review.NewMockrepositoryInterface(ctrl)
	ruc := NewReviewUseCase(mockReviewRepo)
	return ruc, &mocks{
		reviewRepo: mockReviewRepo,
	}
}

func TestUseCaseGetReviews(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		offset  int
		limit   int
		want    []Review
		wantErr bool
	}{}

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
			m.reviewRepo.EXPECT().
				GetReviews(tt.offset, tt.limit).
				Return()

			// テスト対象の実行
			got, err := ruc.GetReviews(tt.offset, tt.limit)
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
