package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
}
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stm := `INSERT INTO users (name, email, hashed_password, created)
          VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stm, name, email, hashedPassword)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicatedEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(name, email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(name, email, password string) (bool, error) {
	return true, nil
}
