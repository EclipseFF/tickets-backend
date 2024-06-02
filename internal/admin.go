package internal

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Admin struct {
	Id       *int     `json:"id"`
	Email    *string  `json:"email"`
	Password Password `json:"-"`
}

type AdminRepo struct {
	DB *pgxpool.Pool
}

func (m *AdminRepo) CreateAdmin(admin *Admin) (*Session, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var id int
	err = tx.QueryRow(context.Background(), `INSERT INTO admin_users(id, email, password) values (default, $1, $2) RETURNING id`, admin.Email, admin.Password.Hash).Scan(&id)
	if err != nil {
		return nil, err
	}
	token, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO admin_session(token, admin_id) values ($1, $2)`, token.String(), id)
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

func (m *AdminRepo) CreateSession(userId *int) (*string, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	token, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO admin_session(token, admin_id) VALUES ($1, $2)`, token.String(), *userId)
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

func (m *AdminRepo) GetAdminByEmail(email *string) (*Admin, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var admin Admin
	err = tx.QueryRow(context.Background(), `SELECT id, email, password FROM admin_users where email = $1`, *email).Scan(&admin.Id, &admin.Email, &admin.Password.Hash)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (m *AdminRepo) GetAdminBySession(token *string) (*Admin, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var id int
	err = tx.QueryRow(context.Background(), `SELECT admin_id FROM admin_session WHERE token = $1`, *token).Scan(&id)
	if err != nil {
		return nil, err
	}
	var admin Admin
	err = tx.QueryRow(context.Background(), `SELECT email FROM admin_users WHERE id = $1`, id).Scan(&admin.Email)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (m *AdminRepo) DeleteSession(token *string) error {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `DELETE FROM admin_session WHERE token = $1`, *token)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	return err
}

func (m *AdminRepo) EnsureSession(token *string) (*int, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var id int
	err = tx.QueryRow(context.Background(), `SELECT admin_id FROM admin_session WHERE token = $1`, *token).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
