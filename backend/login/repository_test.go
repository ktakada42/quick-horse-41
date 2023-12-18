package login

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	userId = "user"
	token  = "token"
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		want        string
		wantErr     bool
	}{
		{
			name: "OK",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT token FROM session WHERE user_id = ?").WithArgs(userId).WillReturnRows(sqlmock.NewRows([]string{"token"}).AddRow(token))
			},
			want:    token,
			wantErr: false,
		},
		{
			name: "OK: no rows returned",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT token FROM session WHERE user_id = ?").WithArgs(userId).WillReturnRows(sqlmock.NewRows([]string{"token"}))
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "NG",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT token FROM session WHERE user_id = ?").WithArgs(userId).WillReturnError(fmt.Errorf("error"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockClosure(mock)

			r := NewLoginRepository(db)

			got, err := r.getToken(userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
