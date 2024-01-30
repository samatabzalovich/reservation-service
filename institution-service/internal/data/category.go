package data

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCategory(id int64, name string, description string) (*Category, error) {
	// validate the input
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
	}, nil
}

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) Insert(category *Category) (int64, error) {
	query := `INSERT INTO category (name, description) VALUES ($1, $2) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int64
	err := m.DB.QueryRowContext(ctx, query, category.Name, category.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m CategoryModel) GetById(id int64) (*Category, error) {
	query := `SELECT id, name, description FROM category WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var category *Category
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (m CategoryModel) GetAll() ([]*Category, error) {
	query := `SELECT id, name, description FROM category`
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
		err = rows.Scan(&category.ID, &category.Name, &category.Description)
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
	query := `UPDATE category SET name = $1, description = $2 WHERE id = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, category.Name, category.Description, category.ID)
	if err != nil {
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
