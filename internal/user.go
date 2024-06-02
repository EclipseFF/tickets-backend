package internal

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRepo struct {
	DB *pgxpool.Pool
}

type User struct {
	Id       *int     `json:"id"`
	Email    *string  `json:"email"`
	Password Password `json:"-"`
	Phone    *string  `json:"phone"`
}

type AdditionalUserData struct {
	UserId      int        `json:"user_id"`
	Surname     *string    `json:"surname"`
	Name        *string    `json:"name"`
	Patronymic  *string    `json:"patronymic"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type Session struct {
	Token  *string `json:"token"`
	UserId *int    `json:"user_id"`
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

func (m *UserRepo) GetUserById(id int) (*User, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	row := tx.QueryRow(context.Background(), `SELECT * FROM users WHERE id = $1`, id)
	var u User
	err = row.Scan(&u.Id, &u.Phone, &u.Email)
	tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *UserRepo) GetUserBySession(token *string) (*User, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var id int
	err = tx.QueryRow(context.Background(), `SELECT user_id FROM sessions WHERE token = $1`, *token).Scan(&id)
	if err != nil {
		return nil, err
	}
	var u User
	err = tx.QueryRow(context.Background(), `SELECT id, email, phone FROM users where id = $1`, id).Scan(&u.Id, &u.Email, &u.Phone)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *UserRepo) CreateUser(user *User) (*Session, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var id int
	err = tx.QueryRow(context.Background(), `INSERT INTO users(id, email, password, phone) values (default, $1, $2, $3) RETURNING id`, user.Email, user.Password.Hash, user.Phone).Scan(&id)
	if err != nil {
		return nil, err
	}
	token, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO sessions(token, user_id) values ($1, $2)`, token.String(), id)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	var s Session
	s.UserId = &id
	temp := token.String()
	s.Token = &temp
	return &s, nil
}

func (m *UserRepo) GetUserByEmail(email *string) (*User, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var u User
	err = tx.QueryRow(context.Background(), `SELECT id, email, password, phone FROM users where email = $1`, email).Scan(&u.Id, &u.Email, &u.Password.Hash, &u.Phone)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *UserRepo) CreateSession(userId *int) (*string, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	token, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO sessions(token, user_id) VALUES ($1, $2)`, token.String(), *userId)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	s := token.String()
	return &s, nil
}

func (m *UserRepo) DeleteSession(token *string) error {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `DELETE FROM sessions WHERE token = $1`, *token)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	return err
}

func (m *UserRepo) GetUserAdditional(id int) (*AdditionalUserData, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var data AdditionalUserData
	err = tx.QueryRow(context.Background(), `SELECT user_id, surname, name, patronymic, date_of_birth FROM additional_user_data WHERE user_id = $1`, id).Scan(&data.UserId, &data.Surname, &data.Name, &data.Patronymic, &data.DateOfBirth)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &data, nil
}
