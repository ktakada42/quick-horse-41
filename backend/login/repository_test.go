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

			got, err := r.GetToken(userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateToken(t *testing.T) {
	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		want        string
		wantErr     bool
	}{
		{
			name: "OK",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO session").WithArgs(userId, token).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "NG",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO session").WithArgs(userId, token).WillReturnError(fmt.Errorf("error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "NG: error at Begin()",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "NG: error at Commit()",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO session").WithArgs(userId, token).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(fmt.Errorf("error"))
			},
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

			if err := r.CreateToken(userId, token); (err != nil) != tt.wantErr {
				t.Fatalf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateToken(t *testing.T) {
	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		want        string
		wantErr     bool
	}{
		{
			name: "OK",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE session").WithArgs(token, userId).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "NG",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE session").WithArgs(token, userId).WillReturnError(fmt.Errorf("error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "NG: error at Begin()",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "NG: error at Commit()",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE session").WithArgs(token, userId).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(fmt.Errorf("error"))
			},
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

			if err := r.UpdateToken(userId, token); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
