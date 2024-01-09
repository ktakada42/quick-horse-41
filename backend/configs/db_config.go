package configs

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DBHost = "mysql"
	DBPort = "3306"
)

var (
	DBUser     = os.Getenv("MYSQL_USER")
	DBPassword = os.Getenv("MYSQL_PASSWORD")
	DBName     = os.Getenv("MYSQL_DATABASE")
)

func NewDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPassword, DBHost, DBPort, DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
