package internal

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *pgxpool.Pool
}
type User struct {
	Id       *int     `json:"id"`
	Username *string  `json:"username"`
	Password Password `json:"-"`
	Email    *string  `json:"email"`
	Status   *string  `json:"status"`
	Role     *string  `json:"role"`
}

type Password struct {
	Plaintext string
	Hash      string
}

func (p *Password) SetPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Plaintext), 14)
	if err != nil {
		return err
	}
	p.Hash = string(hash)
	return nil
}

func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (m *UserModel) GetUser(id int) (*User, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	row := tx.QueryRow(context.Background(), `SELECT id, username, email, status, role_id FROM users WHERE id = $1`, id)
	var u User
	var rId int
	err = row.Scan(&u.Id, &u.Username, &u.Email, &u.Status, &rId)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
