package data

import (
	"authentication-service/internal/validator"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateNumber = errors.New("duplicate number")
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"userName"`
	Number    string    `json:"number"`
	Password  Password  `json:"-"`
	Type      string    `json:"type"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}


func ValidateNumber(v *validator.Validator, number string) {
	v.Check(number != "", "number", "must be provided")
	v.Check(validator.Matches(number, validator.PhoneRX), "number", "must be a valid phone number")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.UserName != "", "name", "must be provided")
	v.Check(len(user.UserName) <= 500, "name", "must not be more than 500 bytes long")

	ValidateNumber(v, user.Number)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}

}

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
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

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username, number, password_hash, activated, type)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, version`
	args := []any{user.UserName, user.Number, user.Password.hash, user.Activated, user.Type}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_number_key" (SQLSTATE 23505)`:
			return ErrDuplicateNumber
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByNumber(number string) (*User, error) {
	query := `
	SELECT id, created_at, username, number, password_hash, activated, version, type
	FROM users
	WHERE number = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, number).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UserName,
		&user.Number,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.Type,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET username = $1, password_hash = $2, activated = $3, version = version + 1 , type = $4
	WHERE id = $5 AND version = $6
	RETURNING version`
	args := []any{
		user.UserName,
		user.Number,
		user.Activated,
		user.Type,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_number_key" (SQLSTATE 23505)`:
			return ErrDuplicateNumber
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (m UserModel) ActivateUser(number string) (int64, error) {
	var id int64
	query := `
	UPDATE users
	SET activated = true, version = version + 1
	WHERE number = $1 AND version = 1
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, number).Scan(&id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrEditConflict
		default:
			return 0, err
		}
	}
	return id, nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
	SELECT users.id, users.created_at, users.username, users.number, users.password_hash, users.activated, users.version, users.type
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UserName,
		&user.Number,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.Type,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
