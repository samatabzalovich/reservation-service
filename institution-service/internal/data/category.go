package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Category struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	PhotoUrl    *string `json:"photoUrl,omitempty"`
}

func NewCategory(id int64, name string, description string, photoUrl string) (*Category, error) {
	// validate the input
	if id < 0 {
		return nil, ErrInvalidCategoryId
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	return &Category{
		ID:          id,
		Name:        name,
		Description: description,
		PhotoUrl:    &photoUrl,
	}, nil
}

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) Insert(category *Category) (int64, error) {
	query := `INSERT INTO category (name, description, photo_url) VALUES ($1, $2, $3) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int64
	err := m.DB.QueryRowContext(ctx, query, category.Name, category.Description, category.PhotoUrl).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m CategoryModel) GetById(id int64) (*Category, error) {
	query := `SELECT id, name, description, photo_url FROM category WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var category Category
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.Description, &category.PhotoUrl)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (m CategoryModel) GetAll() ([]*Category, error) {
	query := `SELECT id, name, description, photo_url FROM category`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*Category
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name, &category.Description, &category.PhotoUrl)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (m CategoryModel) Update(category *Category) error {
	query := `UPDATE category SET name = $1, description = $2, photo_url = $3 WHERE id = $4 RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, category.Name, category.Description, category.PhotoUrl, category.ID).Scan(&category.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (m CategoryModel) Delete(id int64) error {
	query := `DELETE FROM category WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m CategoryModel) GetByInstitution(instId int64) ([]*Category, error) {
	query := `SELECT c.id, c.name, c.description, c.photo_url FROM category c JOIN institution_category ic ON c.id = ic.category_id WHERE ic.inst_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*Category
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name, &category.Description, &category.PhotoUrl)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
