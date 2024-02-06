package data

import (
	"context"
	"database/sql"
	"time"
)

type Service struct {
	ID          int64     `json:"id"`
	InstId      int64     `json:"inst_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Duration    time.Time `json:"duration"`
	PhotoUrl    string    `json:"photo_url"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func NewService(instId int64, name string, description string, price float64, duration time.Time, photoUrl string) (*Service, error) {
	if instId < 1 {
		return nil, ErrInvalidInstId
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if price < 1 {
		return nil, ErrInvalidPrice
	}
	service := &Service{
		InstId:      instId,
		Name:        name,
		Description: description,
		Price:       price,
		Duration:    duration,
		PhotoUrl:    photoUrl,
	}
	return service, nil
}

type ServiceModel struct {
	DB *sql.DB
}

func (m ServiceModel) Insert(service *Service) error {
	query := `
		INSERT INTO services (institution_id, name, description, price, duration, photo_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	args := []interface{}{service.InstId, service.Name, service.Description, service.Price, service.Duration, service.PhotoUrl}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&service.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m ServiceModel) GetAllForInst(instId int64) ([]*Service, error) {
	query := `
		SELECT id, inst_id, name, description, price, duration, photo_url, created_at, updated_at
		FROM services
		WHERE inst_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var services []*Service
	for rows.Next() {
		var service Service
		err := rows.Scan(
			&service.ID,
			&service.InstId,
			&service.Name,
			&service.Description,
			&service.Price,
			&service.Duration,
			&service.PhotoUrl,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func (m ServiceModel) GetById(id int64) (*Service, error) {
	query := `
		SELECT id, inst_id, name, description, price, duration, photo_url, created_at, updated_at
		FROM services
		WHERE id = $1`
	var service Service
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.InstId,
		&service.Name,
		&service.Description,
		&service.Price,
		&service.Duration,
		&service.PhotoUrl,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &service, nil
}

func (m ServiceModel) Update(service *Service) error {
	query := `
		UPDATE services
		SET name = $1, description = $2, price = $3, duration = $4, photo_url = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND inst_id = $7`
	args := []interface{}{
		service.Name,
		service.Description,
		service.Price,
		service.Duration,
		service.PhotoUrl,
		service.ID,
		service.InstId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&service.ID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (m ServiceModel) Delete(id int64) error {
	query := `
		DELETE FROM services
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
