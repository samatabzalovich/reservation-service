package data

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx"
)

type Comment struct {
	ID            int64     `json:"id"`
	Comment       string    `json:"comment"`
	InstitutionId int64     `json:"inst_id"`
	UserId        int64     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewComment(comment string, instId, userId int64) (*Comment, error) {
	if comment == "" {
		return nil, ErrInvalidField
	}
	c := &Comment{
		Comment:       comment,
		InstitutionId: instId,
		UserId:        userId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return c, nil
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(c *Comment) error {
	query := `INSERT INTO comment (comment, institution_id, user_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	args := []interface{}{c.Comment, c.InstitutionId, c.UserId, c.CreatedAt, c.UpdatedAt}
	err := m.DB.QueryRow(query, args...).Scan(&c.ID)
	if err != nil {

		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.Code == "23503" {
				return ErrInvalidField
			}
		}
		return err
	}
	return nil
}

func (m *CommentModel) GetById(id int64) (*Comment, error) {
	query := `SELECT id, comment, inst_id, user_id, created_at, updated_at FROM comments WHERE id = $1`
	c := &Comment{}
	err := m.DB.QueryRow(query, id).Scan(&c.ID, &c.Comment, &c.InstitutionId, &c.UserId, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return c, nil
}

func (m *CommentModel) GetByInstitutionId(instId int64) ([]*Comment, error) {
	query := `SELECT id, comment, inst_id, user_id, created_at, updated_at FROM comments WHERE inst_id = $1`
	rows, err := m.DB.Query(query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	for rows.Next() {
		c := &Comment{}
		err := rows.Scan(&c.ID, &c.Comment, &c.InstitutionId, &c.UserId, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) GetByUserId(userId int64) ([]*Comment, error) {
	query := `SELECT id, comment, inst_id, user_id, created_at, updated_at FROM comments WHERE user_id = $1`
	rows, err := m.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	for rows.Next() {
		c := &Comment{}
		err := rows.Scan(&c.ID, &c.Comment, &c.InstitutionId, &c.UserId, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) Update(c *Comment) error {
	query := `UPDATE comments SET comment = $1, updated_at = $2 WHERE id = $3`
	_, err := m.DB.Exec(query, c.Comment, time.Now(), c.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) Delete(id int64) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) DeleteByInstitutionId(instId int64) error {
	query := `DELETE FROM comments WHERE inst_id = $1`
	_, err := m.DB.Exec(query, instId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) DeleteByUserId(userId int64) error {
	query := `DELETE FROM comments WHERE user_id = $1`
	_, err := m.DB.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) GetAllForInstitution(instId int64) ([]*Comment, error) {
	query := `SELECT id, comment, inst_id, user_id, created_at, updated_at FROM comments WHERE inst_id = $1`
	rows, err := m.DB.Query(query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	for rows.Next() {
		c := &Comment{}
		err := rows.Scan(&c.ID, &c.Comment, &c.InstitutionId, &c.UserId, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) GetAllForUser(userId int64) ([]*Comment, error) {
	query := `SELECT id, comment, inst_id, user_id, created_at, updated_at FROM comments WHERE user_id = $1`
	rows, err := m.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	for rows.Next() {
		c := &Comment{}
		err := rows.Scan(&c.ID, &c.Comment, &c.InstitutionId, &c.UserId, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
