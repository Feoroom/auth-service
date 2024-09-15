package models

import (
	"database/sql"
	"log"
	"time"
)

type Session struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshTokenHash string    `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

type SessionModel struct {
	DB *sql.DB
}

func (sm SessionModel) CreateSession(s *Session) (*Session, error) {
	q := `insert into sessions (id, user_id, refresh_token, expires_at)
	values ($1, $2, $3, $4)`

	args := []interface{}{s.ID, s.UserID, s.RefreshToken, s.ExpiresAt}

	_, err := sm.DB.Exec(q, args...)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (sm SessionModel) GetSession(id string) (*Session, error) {
	var s Session

	q := `select * from sessions
		where user_id=$1`

	log.Print(id)

	err := sm.DB.QueryRow(q, id).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshToken,
		&s.CreatedAt,
		&s.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
