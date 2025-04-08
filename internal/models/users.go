package models

import (
	"database/sql"
	"time"
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
	stm := `INSERT INTO users (name, email, hashed_password, created)
          VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err := m.DB.Exec(stm, name, email)
	if err != nil {
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
