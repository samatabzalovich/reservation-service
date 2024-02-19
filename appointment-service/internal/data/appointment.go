package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)
const dateTimeFormat = "2006-01-02 15:04:05.999999Z"

type Appointment struct {
	ID          int      `json:"id"`
	ClientID    int      `json:"client_id"`
	InstId      int      `json:"inst_id"`
	EmployeeID  int      `json:"employee_id"`
	ServiceID   int      `json:"service_id"`
	StartTime   DateTime `json:"start_time"`
	EndTime     DateTime `json:"end_time"`
	IsCancelled bool     `json:"is_cancelled"`
	CreatedAt   DateTime `json:"created_at"`
	UpdatedAt   DateTime `json:"updated_at"`
}




type DateTime struct {
	time.Time
}

func (ct DateTime) MarshalJSON() ([]byte, error) {
	formattedTime := fmt.Sprintf("\"%s\"", ct.Format(dateTimeFormat))
	return []byte(formattedTime), nil
}

func (ct *DateTime) UnmarshalJSON(data []byte) error {
	// Trim the quotes from the JSON string
	strTime := string(data)
	strTime = strTime[1 : len(strTime)-1] // Remove quotes

	// Parse the time string using the custom format
	parsedTime, err := time.Parse(dateTimeFormat, strTime)
	if err != nil {
		return ErrInvalidAppointmentTime
	}

	ct.Time = parsedTime
	return nil
}

func NewAppointment(
	clientID int,
	instId int,
	employeeID int,
	serviceID int,
	startTime DateTime,
	endTime DateTime,
	isCancelled bool,
) (*Appointment, error) {
	if endTime.Before(startTime.Time) {
		return nil, ErrInvalidAppointmentTime
	}
	if startTime.Local().Before(time.Now()) {
		return nil, ErrInvalidAppointmentTime
	}
	if clientID < 1 {
		return nil, ErrInvalidClientID
	}
	if instId < 1 {
		return nil, ErrInvalidInstID
	}
	if employeeID < 1 {
		return nil, ErrInvalidEmployeeID
	}
	if serviceID < 1 {
		return nil, ErrInvalidServiceID
	}
	return &Appointment{
		ClientID:    clientID,
		InstId:      instId,
		EmployeeID:  employeeID,
		ServiceID:   serviceID,
		StartTime:   DateTime{startTime.Local()},
		EndTime:     DateTime{endTime.Local()},
		IsCancelled: isCancelled,
	}, nil
}

type AppointmentModel struct {
	DB *sql.DB
}

func (m AppointmentModel) Insert(appointment *Appointment) error {
	query := `INSERT INTO appointments (user_id, institution_id, employee_id, service_id, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query,
		appointment.ClientID,
		appointment.InstId,
		appointment.EmployeeID,
		appointment.ServiceID,
		appointment.StartTime.Time,
		appointment.EndTime.Time,
	).Scan(&appointment.ID)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "appointments_employee_id_fkey" {
				return ErrInvalidEmployeeID
			}
			if pgerr.ConstraintName == "appointments_institution_id_fkey" {
				return ErrInvalidInstID
			}
			if pgerr.ConstraintName == "appointments_service_id_fkey" {
				return ErrInvalidServiceID
			}
			if pgerr.ConstraintName == "appointments_user_id_fkey" {
				return ErrInvalidClientID
			}
		}
		return err
	}
	return nil
}

func (m AppointmentModel) GetAllForInst(instId int64) ([]*Appointment, error) {
	query := `SELECT 
		id, user_id, institution_id, employee_id, service_id, start_time, end_time, is_canceled, created_at, updated_at
	 FROM appointments WHERE institution_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		var startTime time.Time
		var endTime time.Time
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.InstId,
			&appointment.EmployeeID,
			&appointment.ServiceID,
			&startTime,
			&endTime,
			&appointment.IsCancelled,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointment.StartTime = DateTime{startTime}
		appointment.EndTime = DateTime{endTime}
		appointment.CreatedAt = DateTime{createdAt}
		appointment.UpdatedAt = DateTime{updatedAt}
		appointments = append(appointments, &appointment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

func (m AppointmentModel) GetById(id int64) (*Appointment, error) {
	query := `SELECT * FROM appointments WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, id)
	var appointment Appointment
	var startTime time.Time
		var endTime time.Time
		var createdAt time.Time
		var updatedAt time.Time
	err := row.Scan(
		&appointment.ID,
		&appointment.ClientID,
		&appointment.InstId,
		&appointment.EmployeeID,
		&appointment.ServiceID,
		&startTime,
		&endTime,
		&appointment.IsCancelled,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	appointment.StartTime = DateTime{startTime}
	appointment.EndTime = DateTime{endTime}
	appointment.CreatedAt = DateTime{createdAt}
	appointment.UpdatedAt = DateTime{updatedAt}
	return &appointment, nil
}

func (m AppointmentModel) Update(appointment *Appointment) error {
	query := `UPDATE appointments SET user_id = $1, institution_id = $2, employee_id = $3, service_id = $4, start_time = $5, end_time = $6, is_cancelled = $7, updated_at = CURRENT_TIMESTAMP WHERE id = $8`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query,
		appointment.ClientID,
		appointment.InstId,
		appointment.EmployeeID,
		appointment.ServiceID,
		appointment.StartTime.Time,
		appointment.EndTime.Time,
		appointment.IsCancelled,
		appointment.ID,
	)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "appointments_employee_id_fkey" {
				return ErrInvalidEmployeeID
			}
			if pgerr.ConstraintName == "appointments_institution_id_fkey" {
				return ErrInvalidInstID
			}
			if pgerr.ConstraintName == "appointments_service_id_fkey" {
				return ErrInvalidServiceID
			}
			if pgerr.ConstraintName == "appointments_user_id_fkey" {
				return ErrInvalidClientID
			}
		}
		return err
	}
	return nil
}

func (m AppointmentModel) Delete(id int64) error {
	query := `DELETE FROM appointments WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m AppointmentModel) GetAllForClient(clientId int64) ([]*Appointment, error) {
	query := `SELECT 
		id, user_id, institution_id, employee_id, service_id, start_time, end_time, is_canceled, created_at, updated_at
	 FROM appointments WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		var startTime time.Time
		var endTime time.Time
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.InstId,
			&appointment.EmployeeID,
			&appointment.ServiceID,
			&startTime,
			&endTime,
			&appointment.IsCancelled,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointment.StartTime = DateTime{startTime}
		appointment.EndTime = DateTime{endTime}
		appointment.CreatedAt = DateTime{createdAt}
		appointment.UpdatedAt = DateTime{updatedAt}
		appointments = append(appointments, &appointment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

func (m AppointmentModel) GetAllForEmployee(employeeId int64) ([]*Appointment, error) {
	query := `SELECT 
		id, user_id, institution_id, employee_id, service_id, start_time, end_time, is_canceled, created_at, updated_at
	 FROM appointments WHERE employee_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		var startTime time.Time
		var endTime time.Time
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.InstId,
			&appointment.EmployeeID,
			&appointment.ServiceID,
			&startTime,
			&endTime,
			&appointment.IsCancelled,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointment.StartTime = DateTime{startTime}
		appointment.EndTime = DateTime{endTime}
		appointment.CreatedAt = DateTime{createdAt}
		appointment.UpdatedAt = DateTime{updatedAt}
		appointments = append(appointments, &appointment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}
