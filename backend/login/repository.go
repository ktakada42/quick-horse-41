package login

import (
	"database/sql"
	"errors"
)

type repositoryInterface interface {
	GetToken(userId string) (string, error)
	SaveToken(userId, token string) error
	UpdateToken(userId, token string) error
}

type repositoryStruct struct {
	db *sql.DB
}

func NewLoginRepository(db *sql.DB) repositoryInterface {
	return &repositoryStruct{db: db}
}

func (r *repositoryStruct) GetToken(userId string) (string, error) {
	row := r.db.QueryRow("SELECT token FROM session WHERE user_id = ?", userId)

	var token string
	if err := row.Scan(&token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return token, nil
}

func (r *repositoryStruct) SaveToken(userId, token string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO session (user_id, token) VALUES (?, ?)", userId, token)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *repositoryStruct) UpdateToken(userId, token string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE session SET token = ? WHERE user_id = ?", token, userId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
