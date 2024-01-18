package controller

import (
	"app/review/entity"
	"app/review/mock/mock_usecase"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	reviewUseCase *mock_usecase.MockUseCaseInterface
}

func newWithMocks(t *testing.T) (ControllerInterface, *mocks) {
	ctrl := gomock.NewController(t)
	mockReviewUseCase := mock_usecase.NewMockUseCaseInterface(ctrl)
	rc := NewReviewController(mockReviewUseCase)
	return rc, &mocks{
		reviewUseCase: mockReviewUseCase,
	}
}

func TestGetReviews(t *testing.T) {
	t.Parallel()

	// テストデータの準備
	reviews := setTestData()

	type mock struct {
		res []entity.Review
		err error
	}

	tests := map[string]struct {
		requestBody    string
		mock           mock
		expectedStatus int
	}{
		"正常系: offset = 0, limit = 10": {
			requestBody: `{"offset": 0, "limit": 10}`,
			mock: mock{
				res: reviews,
				err: nil,
			},
			expectedStatus: http.StatusOK,
		},
		"正常系: offset がデフォルト値": {
			requestBody: `{"limit": 10}`,
			mock: mock{
				res: reviews,
				err: nil,
			},
			expectedStatus: http.StatusOK,
		},
		"異常系: offset が負の値": {
			requestBody:    `{"offset": -1, "limit": 10}`,
			expectedStatus: http.StatusBadRequest,
		},
		"異常系: limit が 1 より小さい": {
			requestBody:    `{"offset": 0, "limit": -1}`,
			expectedStatus: http.StatusBadRequest,
		},
		"異常系: limit が 100 より大きい": {
			requestBody:    `{"offset": 0, "limit": 1000}`,
			expectedStatus: http.StatusBadRequest,
		},
		"異常系: offset と limit が範囲外": {
			requestBody:    `{"offset": -1, "limit": -1}`,
			expectedStatus: http.StatusBadRequest,
		},
		"異常系: internal server error": {
			requestBody: `{"offset": 0, "limit": 10}`,
			mock: mock{
				res: nil,
				err: errors.New("internal server error"),
			},
			expectedStatus: http.StatusInternalServerError,
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
			rc, m := newWithMocks(t)

			// モックの設定
			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusInternalServerError {
				m.reviewUseCase.EXPECT().GetReviews(gomock.Any(), gomock.Any()).
					Return(tt.mock.res, tt.mock.err)
			}

			// HTTPリクエストとレスポンスの設定
			req, err := http.NewRequest("GET", "/reviews", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Errorf("http.NewRequest() error: %v", err)
			}
			responseRecorder := httptest.NewRecorder()

			// テスト対象の実行
			rc.GetReviews(responseRecorder, req)

			// 結果の検証
			assert.Equal(t, tt.expectedStatus, responseRecorder.Code)
		})
	}
}
