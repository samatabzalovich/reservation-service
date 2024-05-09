package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

type Service struct {
	ID          int64         `json:"id"`
	InstId      int64         `json:"inst_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       int           `json:"price"`
	Duration    time.Duration `json:"duration"`
	PhotoUrl    string        `json:"photo_url"`
	CreatedAt   string        `json:"created_at"`
	ServiceType string        `json:"serviceType"`
	UpdatedAt   string        `json:"updated_at"`
}

func NewService(instId int64, name string, description string, price int, duration time.Duration, photoUrl, serviceType string) (*Service, error) {
	existingServiceTypes := []string{"appointment", "walk-in"}
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
	if serviceType == "" {
		return nil, ErrInvalidServiceType
	}

	if !contains(existingServiceTypes, serviceType) {
		return nil, ErrInvalidServiceType

	}
	service := &Service{
		InstId:      instId,
		Name:        name,
		Description: description,
		Price:       price,
		Duration:    duration,
		PhotoUrl:    photoUrl,
		ServiceType: serviceType,
	}
	return service, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false

}

type ServiceModel struct {
	DB *sql.DB
}

func (m ServiceModel) Insert(service *Service) error {
	intervalStr := fmt.Sprintf("%d hours %d minutes %d seconds",
		int(service.Duration.Hours()), int(service.Duration.Minutes())%60, int(service.Duration.Seconds())%60)
	query := `
		INSERT INTO services (institution_id, name, description, price, duration, photo_url, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	args := []interface{}{service.InstId, service.Name, service.Description, service.Price, intervalStr, service.PhotoUrl, service.ServiceType}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	log.Println(args)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&service.ID)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.Code == "23503" {
				return ErrInvalidInstId
			}
		}
		return err
	}
	return nil
}

func (m ServiceModel) GetAllForInst(instId int64) ([]*Service, error) {
	query := `
		SELECT id, institution_id, name, description, price, duration, photo_url, created_at, updated_at, type
		FROM services
		WHERE institution_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var services []*Service
	for rows.Next() {
		var durationStr string
		var service Service
		err := rows.Scan(
			&service.ID,
			&service.InstId,
			&service.Name,
			&service.Description,
			&service.Price,
			&durationStr,
			&service.PhotoUrl,
			&service.CreatedAt,
			&service.UpdatedAt,
			&service.ServiceType,
		)
		if err != nil {
			return nil, err
		}
		// duration 05:00:00 TO 5h0m0s
		durationStr = strings.Replace(durationStr, ":", "h", 1)
		durationStr = strings.Replace(durationStr, ":", "m", 1)
		parsedDuration, err := time.ParseDuration(durationStr + "s")
		if err != nil {
			return nil, err
		}
		service.Duration = parsedDuration
		services = append(services, &service)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func (m ServiceModel) GetById(id int64) (*Service, error) {
	query := `
		SELECT id, institution_id, name, description, price, duration, photo_url, created_at, updated_at, type
		FROM services
		WHERE id = $1`
	var service Service
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var durationStr string
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.InstId,
		&service.Name,
		&service.Description,
		&service.Price,
		&durationStr,
		&service.PhotoUrl,
		&service.CreatedAt,
		&service.UpdatedAt,
		&service.ServiceType,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	durationStr = strings.Replace(durationStr, ":", "h", 1)
	durationStr = strings.Replace(durationStr, ":", "m", 1)
	parsedDuration, err := time.ParseDuration(durationStr + "s")
	if err != nil {
		return nil, err
	}
	service.Duration = parsedDuration
	return &service, nil
}

func (m ServiceModel) Update(service *Service) error {
	query := `
		UPDATE services
		SET name = $1, description = $2, price = $3, duration = $4, photo_url = $5, updated_at = CURRENT_TIMESTAMP, type = $6
		WHERE id = $7 AND institution_id = $8`
	args := []interface{}{
		service.Name,
		service.Description,
		service.Price,
		service.Duration,
		service.PhotoUrl,
		service.ServiceType,
		service.ID,
		service.InstId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&service.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
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
