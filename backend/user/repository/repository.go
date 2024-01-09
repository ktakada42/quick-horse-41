package repository

import (
	"database/sql"
	"errors"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type UserRepository interface {
	GetHashedPassword(userId string) (string, error)
	GetUserIdByEmail(email string) (string, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetHashedPassword(userId string) (string, error) {
	row := r.db.QueryRow("SELECT password FROM user WHERE user_id = ?", userId)

	var hashedPassword string
	if err := row.Scan(&hashedPassword); err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func (r *userRepository) GetUserIdByEmail(email string) (string, error) {
	row := r.db.QueryRow("SELECT user_id FROM user WHERE email = ?", email)

	var userId string
	if err := row.Scan(&userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return userId, nil
}
