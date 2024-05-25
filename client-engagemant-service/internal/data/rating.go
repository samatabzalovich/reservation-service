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
	ID            int64     `json:"id"`
	AppointmentId int64     `json:"appointment_id"`
	EmployeeId    int64     `json:"employee_id"`
	ClientId      int64     `json:"client_id"`
	InstitutionId int64     `json:"institution_id"`
	ServiceId     int64     `json:"service_id"`
	Rating        int       `json:"rating"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"created_at"`
	UpdateAt      time.Time `json:"updated_at"`
}

type IdForRating struct {
	EmployeeId    int64 `json:"employee_id"`
	InstitutionId int64 `json:"institution_id"`
	ServiceId     int64 `json:"service_id"`
}

func NewRating(appointmentId, employeeId, clientId, institutionId int64, rating int, comment string, serviceId int64) (*Rating, error) {
	if rating < 1 || rating > 10 {
		return nil, ErrRatingMustBe1
	}
	if appointmentId < 1 {
		return nil, ErrInvalidAppointmentId
	}
	if employeeId < 1 {
		return nil, ErrInvalidEmployeeId

	}
	if clientId < 1 {
		return nil, ErrInvalidClientId

	}
	if institutionId < 1 {
		return nil, ErrInvalidInstitutionId

	}
	if serviceId < 1 {
		return nil, ErrInvalidServiceId

	}
	r := &Rating{
		AppointmentId: appointmentId,
		EmployeeId:    employeeId,
		ClientId:      clientId,
		InstitutionId: institutionId,
		Rating:        rating,
		Comment:       comment,
		ServiceId:     serviceId,
	}
	return r, nil
}

type RatingModel struct {
	DB *sql.DB
}

func (m *RatingModel) Insert(rating *Rating) error {
	query := `
		INSERT INTO rating (appointment_id, employee_id, user_id, institution_id, rating, comment, service_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, rating.AppointmentId, rating.EmployeeId, rating.ClientId, rating.InstitutionId, rating.Rating, rating.Comment, rating.ServiceId).Scan(&rating.ID, &rating.CreatedAt, &rating.UpdateAt)
	if err != nil {

		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.Code == "23503" {
				return ErrInvalidField
			}
			if pgerr.Code == "23505" {
				return ErrRatingAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (m *RatingModel) GetServiceIdAndInstIdToRatingAppointment(appointmentId, userId int64) (*IdForRating, error) {
	query := `
		SELECT service_id, institution_id, employee_id
		FROM appointments
		WHERE id = $1 AND user_id = $2
	`
	var ids IdForRating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, appointmentId, userId).Scan(&ids.ServiceId, &ids.InstitutionId, &ids.EmployeeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidAppointmentId
		}
		return nil, err
	}
	return &ids, nil
}

func (m *RatingModel) GetAllForInst(instId int64) ([]*Rating, error) {
	query := `
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, service_id, created_at, updated_at
		FROM rating
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
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.ServiceId, &r.CreatedAt, &r.UpdateAt)
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
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, created_at, updated_at, service_id
		FROM rating
		WHERE id = $1
	`
	var rating Rating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&rating.ID, &rating.AppointmentId, &rating.EmployeeId, &rating.ClientId, &rating.InstitutionId, &rating.Rating, &rating.Comment, &rating.CreatedAt, &rating.UpdateAt, &rating.ServiceId)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (m *RatingModel) Update(rating *Rating) error {
	query := `
		UPDATE rating
		SET rating = $1, comment = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, rating.Rating, rating.Comment, rating.ID, rating.ClientId).Scan(&rating.UpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.Code == "23514" {
				return ErrRatingMustBe1
			}
		}
		return err
	}
	return nil
}

func (m *RatingModel) Delete(id int64) error {
	query := `
		DELETE FROM rating
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
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, created_at, updated_at, service_id
		FROM rating
		WHERE user_id = $1
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
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt, &r.ServiceId)
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
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, created_at, updated_at, service_id
		FROM rating
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
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt, &r.ServiceId)
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
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, created_at, updated_at, service_id
		FROM rating
		WHERE appointment_id = $1
	`
	var rating Rating
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, appointmentId).Scan(&rating.ID, &rating.AppointmentId, &rating.EmployeeId, &rating.ClientId, &rating.InstitutionId, &rating.Rating, &rating.Comment, &rating.CreatedAt, &rating.UpdateAt, &rating.ServiceId)
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
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, created_at, updated_at, service_id
		FROM rating
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
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt, &r.ServiceId)
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
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, created_at, updated_at, service_id
		FROM rating
		WHERE user_id = $1
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
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt, &r.ServiceId)
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
		SELECT id, appointment_id, employee_id, user_id, institution_id, rating, comment, created_at, updated_at, service_id
		FROM rating
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
		err := rows.Scan(&r.ID, &r.AppointmentId, &r.EmployeeId, &r.ClientId, &r.InstitutionId, &r.Rating, &r.Comment, &r.CreatedAt, &r.UpdateAt, &r.ServiceId)
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
		FROM rating
		WHERE employee_id = $1
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, employeeId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}

func (m *RatingModel) GetAverageRatingForInstitution(instId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM rating
		WHERE institution_id = $1
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, instId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}

func (m *RatingModel) GetAverageRatingForService(serviceId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM rating
		WHERE service_id = $1
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, serviceId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}

func (m *RatingModel) GetAverageRatingForClient(clientId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM rating
		WHERE user_id = $1
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}

func (m *RatingModel) GetAverageRatingForEmployeeService(employeeId, serviceId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM rating
		WHERE employee_id = $1 AND service_id = $2
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, employeeId, serviceId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}

func (m *RatingModel) GetAverageRatingForClientService(clientId, serviceId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM rating
		WHERE user_id = $1 AND service_id = $2
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId, serviceId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}

func (m *RatingModel) GetAverageRatingForClientInstitution(clientId, instId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM rating
		WHERE user_id = $1 AND institution_id = $2
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId, instId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}

func (m *RatingModel) GetAverageRatingForClientEmployee(clientId, employeeId int64) (float64, error) {
	query := `
		SELECT AVG(rating)
		FROM rating
		WHERE user_id = $1 AND employee_id = $2
	`
	var avgRating *float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, clientId, employeeId).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	if avgRating == nil {
		return 0, nil
	}
	return *avgRating, nil
}
