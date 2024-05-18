package data

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx"
)

const (
	TimeParse      = "15:04:00"
	dateTimeFormat = time.RFC3339
)

type Appointment struct {
	ID          int64    `json:"id"`
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

type EmployeeSchedule struct {
	DayOfWeek      int       `json:"day_of_week"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	BreakStartTime time.Time `json:"break_start_time"`
	BreakEndTime   time.Time `json:"break_end_time"`
}

type SlotResponse struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (es *EmployeeSchedule) ParseToTimeStruct(dayOfWeek int, startTime string, endTime string, breakStartTime string, breakEndTime string) error {
	var err error
	es.StartTime, err = time.Parse(TimeParse, startTime)
	if err != nil {
		return err
	}
	es.EndTime, err = time.Parse(TimeParse, endTime)
	if err != nil {
		return err
	}
	es.BreakStartTime, err = time.Parse(TimeParse, breakStartTime)
	if err != nil {
		return err
	}
	es.BreakEndTime, err = time.Parse(TimeParse, breakEndTime)
	if err != nil {
		return err
	}
	es.DayOfWeek = dayOfWeek
	return nil
}

type DateTime struct {
	time.Time
}

// func (ct DateTime) MarshalJSON() ([]byte, error) {
// 	formattedTime := fmt.Sprintf("\"%s\"", ct.Format(dateTimeFormat))
// 	return []byte(formattedTime), nil
// }

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
	availabeTimeSlots []*SlotResponse,
) (*Appointment, error) {
	startTime.Time = startTime.Time.Local()
	endTime.Time = endTime.Time.Local()
	if availabeTimeSlots == nil {
		return nil, ErrInvalidAppointment
	}
	if startTime.Time.Equal(endTime.Time) {
		return nil, ErrInvalidAppointmentTime
	}
	if endTime.Before(startTime.Time) {
		return nil, ErrInvalidAppointmentTime
	}
	if startTime.Local().Before(time.Now()) {
		return nil, ErrInvalidAppointmentTime
	}

	// check if the selected time slot is available
	isAvailable := false
	for _, slot := range availabeTimeSlots {
		startHour, startMinute, startSecond := startTime.Time.Clock()
		endHour, endMinute, endSecond := endTime.Time.Clock()
		slotStartHour, slotStartMinute, slotStartSecond := slot.StartTime.Clock()
		if startHour == slotStartHour && startMinute == slotStartMinute && startSecond == slotStartSecond && endHour == slot.EndTime.Hour() && endMinute == slot.EndTime.Minute() && endSecond == slot.EndTime.Second() {
			isAvailable = true
			break
		}
	}

	log.Println("isAvailable: ", isAvailable)

	if !isAvailable {
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
		StartTime:   DateTime{startTime.Time},
		EndTime:     DateTime{endTime.Time},
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
	query := `SELECT id, user_id, institution_id, employee_id, service_id, start_time, end_time, is_canceled, created_at, updated_at FROM appointments WHERE id = $1`
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
	query := `UPDATE appointments SET start_time = $1, end_time = $2, is_canceled = $3, updated_at = NOW() WHERE id = $4`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	log.Println("appointment updating: ", appointment.IsCancelled)
	_, err := m.DB.ExecContext(ctx, query,
		appointment.StartTime.Time,
		appointment.EndTime.Time,
		appointment.IsCancelled,
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

func (m AppointmentModel) GetAvailableTimeSlots(instId int64, employeeId int64, serviceId int64, selectedDay DateTime) ([]DateTime, error) {
	query := `SELECT start_time FROM appointments WHERE institution_id = $1 AND employee_id = $2 AND service_id = $3 AND DATE(start_time) = $4`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId, employeeId, serviceId, selectedDay.Time)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var timeSlots []DateTime
	for rows.Next() {
		var startTime time.Time
		err := rows.Scan(&startTime)
		if err != nil {
			return nil, err
		}
		timeSlots = append(timeSlots, DateTime{startTime})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return timeSlots, nil
}

func (m AppointmentModel) GetNumberOfCompletedAppointmentsForUser(userId, instId, employeeId int64) (int, error) {
	query := `SELECT COUNT(id) FROM appointments WHERE user_id = $1 AND end_time < NOW() AND is_canceled = false AND institution_id = $2`
	if employeeId > 0 {
		query += ` AND employee_id = $3`
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, userId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}