package data

import (
	"context"
	"database/sql"
	"time"
)

type Appointment struct {
	ID          int       `json:"id"`
	ClientID    int       `json:"client_id"`
	InstId      int       `json:"inst_id"`
	EmployeeID  int       `json:"employee_id"`
	ServiceID   int       `json:"service_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsCancelled bool      `json:"is_cancelled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewAppointment(
	clientID int,
	instId int,
	employeeID int,
	serviceID int,
	startTime time.Time,
	endTime time.Time,
	isCancelled bool,
) (*Appointment, error) {
	if endTime.Before(startTime) {
		return nil, ErrInvalidAppointmentTime
	}
	if startTime.Before(time.Now().Add(6 * time.Hour)) {
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
		StartTime:   startTime,
		EndTime:     endTime,
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
		appointment.StartTime,
		appointment.EndTime,
	).Scan(&appointment.ID)
	if err != nil {
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
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.InstId,
			&appointment.EmployeeID,
			&appointment.ServiceID,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.IsCancelled,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
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
	err := row.Scan(
		&appointment.ID,
		&appointment.ClientID,
		&appointment.InstId,
		&appointment.EmployeeID,
		&appointment.ServiceID,
		&appointment.StartTime,
		&appointment.EndTime,
		&appointment.IsCancelled,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (m AppointmentModel) Update(appointment *Appointment) error {
	query := `UPDATE appointments SET user_id = $1, institution_id = $2, employee_id = $3, service_id = $4, start_time = $5, end_time = $6, is_cancelled = $7, updated_at = $8 WHERE id = $9`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query,
		appointment.ClientID,
		appointment.InstId,
		appointment.EmployeeID,
		appointment.ServiceID,
		appointment.StartTime,
		appointment.EndTime,
		appointment.IsCancelled,
		appointment.UpdatedAt,
		appointment.ID,
	)
	if err != nil {
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
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.InstId,
			&appointment.EmployeeID,
			&appointment.ServiceID,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.IsCancelled,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
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
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.InstId,
			&appointment.EmployeeID,
			&appointment.ServiceID,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.IsCancelled,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}
