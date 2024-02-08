package data

import (
	"authentication-service/internal/validator"
	"context" // New import
	"crypto/rand"
	"crypto/sha256"
	"database/sql" // New import
	"encoding/base32"
	"log"
	"time"

	"github.com/jackc/pgx"
)

const (
	ScopeEmployeeReg    = "employee_registration"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext     string    `json:"token"`
	Hash          []byte    `json:"-"`
	UserID        int64     `json:"-"`
	Expiry        time.Time `json:"expiry"`
	InstitutionId int64     `json:"-"`
	Scope         string    `json:"-"`
}

func generateToken(userID int64, ttl time.Duration, scope string, instId int64) (*Token, error) {
	token := &Token{
		UserID:        userID,
		Expiry:        time.Now().Add(ttl),
		Scope:         scope,
		InstitutionId: instId,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string, instId int64) (*Token, error) {
	token, err := generateToken(userID, ttl, scope, instId)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	return token, err
}

// Insert() adds the data for a specific token to the tokens table.
func (m TokenModel) Insert(token *Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope, institution_id)
	VALUES ($1, $2, $3, $4, $5)`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}
	if token.InstitutionId != 0 {
		args = append(args, token.InstitutionId)
	} else {
		args = append(args, nil)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			log.Println(pgerr.Code)
			if pgerr.Code == "23503" {
				return ErrRecordNotFound
			}
		} 
	}
	return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}

type MockTokenModel struct {
	DB *sql.DB
}

func (m MockTokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	return nil, nil
}

func (m MockTokenModel) Insert(token *Token) error {
	return nil
}

func (m MockTokenModel) DeleteAllForUser(scope string, userID int64) error {
	return nil
}
