package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx"
)

const (
	TimeParse      = "15:04:00"
	dateTimeFormat = time.RFC3339
)

type Rating struct {
	ID        int64     `json:"id"`
	AppointmentId int64 `json:"appointment_id"`
	EmployeeId int64 `json:"employee_id"`
	ClientId int64 `json:"client_id"`
	InstitutionId int64 `json:"institution_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}



func NewRating(appointmentId, employeeId, clientId, institutionId int64, rating int, comment string) (*Rating, error) {

	r := &Rating{
		AppointmentId: appointmentId,
		EmployeeId: employeeId,
		ClientId: clientId,
		InstitutionId: institutionId,
		Rating: rating,
		Comment: comment,
	}
	return r, nil
}
	

type RatingModel struct {
	DB *sql.DB
}

func (m *RatingModel) Insert(rating *Rating) error {
	query := `
		INSERT INTO ratings (appointment_id, employee_id, client_id, institution_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, rating.AppointmentId, rating.EmployeeId, rating.ClientId, rating.InstitutionId, rating.Rating, rating.Comment).Scan(&rating.ID, &rating.CreatedAt, &rating.UpdateAt)
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

func (m *RatingModel) GetAllForInst(instId int64) ([]*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE institution_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ratings []*Rating
	for rows.Next() {
		r := &Rating{}
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (m *RatingModel) GetById(id int64) (*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE id = $1
	`
	var rating Rating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&rating.ID, &rating.AppointmentId, &rating.EmployeeId, &rating.ClientId, &rating.InstitutionId, &rating.Rating, &rating.Comment, &rating.CreatedAt, &rating.UpdateAt)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (m *RatingModel) Update(rating *Rating) error {
	query := `
		UPDATE ratings
		SET rating = $1, comment = $2, updated_at = NOW()
		WHERE id = $4 AND client_id = $3
		RETURNING updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, rating.Rating, rating.Comment, time.Now(), rating.ID, rating.ClientId).Scan(&rating.UpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (m *RatingModel) Delete(id int64) error {
	query := `
		DELETE FROM ratings
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *RatingModel) GetAllForClient(clientId int64) ([]*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE client_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ratings []*Rating
	for rows.Next() {
		r := &Rating{}
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (m *RatingModel) GetAllForEmployee(employeeId int64) ([]*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE employee_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ratings []*Rating
	for rows.Next() {
		r := &Rating{}
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ratings, nil
}


func (m *RatingModel) GetRatingForAppointment(appointmentId int64) (*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE appointment_id = $1
	`
	var rating Rating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, appointmentId).Scan(&rating.ID, &rating.AppointmentId, &rating.EmployeeId, &rating.ClientId, &rating.InstitutionId, &rating.Rating, &rating.Comment, &rating.CreatedAt, &rating.UpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &rating, nil
}

func (m *RatingModel) GetRatingsForEmployee(employeeId int64) ([]*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE employee_id = $1
	`
	var ratings []*Rating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r := &Rating{}
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (m *RatingModel) GetRatingsForClient(clientId int64) ([]*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE client_id = $1
	`
	var ratings []*Rating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r := &Rating{}
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (m *RatingModel) GetRatingsForInstitution(instId int64) ([]*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, client_id, institution_id, rating, comment, created_at, updated_at
		FROM ratings
		WHERE institution_id = $1
	`
	var ratings []*Rating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r := &Rating{}
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (m *RatingModel) GetAverageRatingForEmployee(employeeId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE employee_id = $1
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, employeeId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}

func (m *RatingModel) GetAverageRatingForInstitution(instId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE institution_id = $1
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, instId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}

func (m *RatingModel) GetAverageRatingForService(serviceId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE service_id = $1
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, serviceId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}

func (m *RatingModel) GetAverageRatingForClient(clientId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE client_id = $1
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}

func (m *RatingModel) GetAverageRatingForAppointment(appointmentId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE appointment_id = $1
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, appointmentId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}



func (m *RatingModel) GetAverageRatingForEmployeeService(employeeId, serviceId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE employee_id = $1 AND service_id = $2
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, employeeId, serviceId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}

func (m *RatingModel) GetAverageRatingForClientService(clientId, serviceId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE client_id = $1 AND service_id = $2
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId, serviceId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}

func (m *RatingModel) GetAverageRatingForClientInstitution(clientId, instId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE client_id = $1 AND institution_id = $2
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId, instId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}


func (m *RatingModel) GetAverageRatingForClientEmployee(clientId, employeeId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM ratings
		WHERE client_id = $1 AND employee_id = $2
	`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId, employeeId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}
