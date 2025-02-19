package data

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx"
)

type DeviceToken struct {
	ID        int64  `json:"id,omitempty"`
	UserID    int64  `json:"user_id"`
	Token     string `json:"token"`
	CreatedAt time.Time
}

type DeviceTokenModel struct {
	DB *sql.DB
}

func (m DeviceTokenModel) Insert(token DeviceToken) (int64, error) {
	stmt := `INSERT INTO user_devices (user_id, token)
	VALUES($1, $2) RETURNING id`
	var id int64
	err := m.DB.QueryRow(stmt, token.UserID, token.Token).Scan(&id)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "user_devices_user_id_fkey" {
				return 0, ErrUserNotFound
			}
			if pgerr.ConstraintName == "user_devices_token_key" {
				return 0, ErrTokenAlreadyExists
			}

		}
		return 0, err
	}
	return id, nil
}

func (m DeviceTokenModel) GetByToken(token string) (*DeviceToken, error) {
	stmt := `SELECT id, user_id, token, created_at FROM user_devices WHERE token = $1`
	row := m.DB.QueryRow(stmt, token)
	var t DeviceToken
	err := row.Scan(&t.ID, &t.UserID, &t.Token, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (m DeviceTokenModel) GetByUserID(userID int64) ([]*DeviceToken, error) {
	stmt := `SELECT id, user_id, token, created_at FROM user_devices WHERE user_id = $1`
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*DeviceToken
	for rows.Next() {
		var t DeviceToken
		err := rows.Scan(&t.ID, &t.UserID, &t.Token, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (m DeviceTokenModel) GetByEmployeeID(employeeId int64) ([]*DeviceToken, error) {
	stmt := `SELECT d.id, d.user_id, d.token   FROM user_devices d join employee e on d.user_id = e.user_id WHERE e.id = $1`
	rows, err := m.DB.Query(stmt, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*DeviceToken
	for rows.Next() {
		var t DeviceToken
		err := rows.Scan(&t.ID, &t.UserID, &t.Token)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (m DeviceTokenModel) Update(token DeviceToken) error {
	stmt := `UPDATE user_devices SET token = $1, updated_at = NOW() WHERE id = $2`
	_, err := m.DB.Exec(stmt, token.Token, token.ID)
	return err
}

func (m DeviceTokenModel) Delete(id int64) error {
	stmt := `DELETE FROM user_devices WHERE id = $1`
	_, err := m.DB.Exec(stmt, id)
	return err
}

func (m DeviceTokenModel) DeleteByToken(token string) error {
	stmt := `DELETE FROM user_devices WHERE token = $1`
	_, err := m.DB.Exec(stmt, token)
	return err
}
