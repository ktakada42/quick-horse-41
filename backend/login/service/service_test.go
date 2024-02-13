package service

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"app/login/mock/mock_repository"
	mock_user_repository "app/user/mock/mock_repository"
)

type loginServiceTest struct {
	lr *mock_repository.MockLoginRepository
	ur *mock_user_repository.MockUserRepository
	ls LoginService
}

func newLoginServiceTest(t *testing.T) *loginServiceTest {
	t.Helper()

	ctrl := gomock.NewController(t)
	lr := mock_repository.NewMockLoginRepository(ctrl)
	ur := mock_user_repository.NewMockUserRepository(ctrl)

	return &loginServiceTest{lr: lr, ur: ur, ls: NewLoginService(lr, ur)}
}

func TestIsPasswordCorrect(t *testing.T) {
	t.Parallel()

	const (
		email    = "test@example.com"
		password = "password"
	)

	hp, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		t.Fatalf("Unexpected error '%v' occurred when generating hashed password", err)
	}
	hashedPassword := string(hp)

	tests := []struct {
		name    string
		expects func(lst *loginServiceTest)
		want    bool
		wantErr bool
	}{
		{
			name: "OK",
			expects: func(lst *loginServiceTest) {
				lst.ur.EXPECT().GetHashedPassword(email).Return(hashedPassword, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "OK: password is incorrect",
			expects: func(lst *loginServiceTest) {
				lst.ur.EXPECT().GetHashedPassword(email).Return("wrongPassword", nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "NG: no password found",
			expects: func(lst *loginServiceTest) {
				lst.ur.EXPECT().GetHashedPassword(email).Return("", fmt.Errorf("error"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lst := newLoginServiceTest(t)
			tt.expects(lst)

			got, err := lst.ls.IsPasswordCorrect(email, password)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPasswordCorrect() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
