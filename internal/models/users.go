package models

import (
	"database/sql"
	"github.com/google/uuid"
	"net"
	"net/http"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) CreateUser(user *User) error {
	q := `insert into users(username, email)
		values ($1, $2)
		returning id`

	err := u.DB.QueryRow(q, user.Username, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil

}

func (u UserModel) GetUser(id string) (*User, error) {
	q := `select * from users
		where id=$1`

	var user User

	err := u.DB.QueryRow(q, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP := net.ParseIP(ip).String()
	return netIP, nil
}
