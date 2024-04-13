package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

var (
	TimeParse = "15:04:00"
)

type Employee struct {
	ID          int64               `json:"id"`
	CreatedAt   time.Time           `json:"created_at"`
	InstId      int64               `json:"inst_id"`
	UserId      int64               `json:"user_id"`
	Name        string              `json:"name"`
	PhotoUrl    string              `json:"photo_url"`
	Description string              `json:"description"`
	Schedule    []*EmployeeSchedule `json:"schedule"`
	Services    []*EmployeeServices `json:"services"`
}

type EmployeeSchedule struct {
	DayOfWeek      int       `json:"day_of_week"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	BreakStartTime time.Time `json:"break_start_time"`
	BreakEndTime   time.Time `json:"break_end_time"`
}

type EmployeeServices struct {
	ServiceId int64  `json:"service_id"`
	Name      string `json:"name"`
}

// Auxiliary types
type employeeScheduleAux struct {
	DayOfWeek      int    `json:"day_of_week"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	BreakStartTime string `json:"break_start_time"`
	BreakEndTime   string `json:"break_end_time"`
}

type TypeForEmployeeTimeSlots struct {
	Schedule ScheduleForEmployeeTimeSlots `json:"schedule"`
	Service  ServiceForEmployeeTimeSlots         `json:"service"`
}

type ScheduleForEmployeeTimeSlots struct {
	DayOfWeek      int       `json:"day_of_week"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	BreakStartTime string `json:"break_start_time"`
	BreakEndTime   string `json:"break_end_time"`
}

type ServiceForEmployeeTimeSlots struct {
	Name     string        `json:"name"`
	Duration string `json:"duration"`
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

func (es *EmployeeSchedule) UnmarshalJSON(data []byte) error {
	var aux employeeScheduleAux
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	if es.StartTime, err = time.Parse(TimeParse, aux.StartTime); err != nil {
		return err
	}
	if es.EndTime, err = time.Parse(TimeParse, aux.EndTime); err != nil {
		return err
	}
	if es.BreakStartTime, err = time.Parse(TimeParse, aux.BreakStartTime); err != nil {
		return err
	}
	if es.BreakEndTime, err = time.Parse(TimeParse, aux.BreakEndTime); err != nil {
		return err
	}

	es.DayOfWeek = aux.DayOfWeek
	return nil
}

func NewEmployee(instId int64, userId int64, description string, schedule []*EmployeeSchedule, services []*EmployeeServices) (*Employee, error) {
	if instId < 1 {
		return nil, ErrInvalidInstId
	}
	if userId < 1 {
		return nil, ErrInvalidUserId
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if len(schedule) < 1 {
		return nil, ErrInvalidSchedule
	}
	if len(services) < 1 {
		return nil, ErrInvalidServices
	}
	employee := &Employee{
		InstId:      instId,
		UserId:      userId,
		Description: description,
		Schedule:    schedule,
		Services:    services,
	}
	return employee, nil
}
func NewEmployeeServices(serviceId int64) (*EmployeeServices, error) {
	if serviceId < 1 {
		return nil, ErrInvalidServiceId // Assuming ErrInvalidServiceId is a predefined error
	}

	return &EmployeeServices{
		ServiceId: serviceId,
	}, nil
}
func NewEmployeeSchedule(dayOfWeek int, startTime time.Time, endTime time.Time, breakStartTime time.Time, breakEndTime time.Time) (*EmployeeSchedule, error) {
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return nil, ErrInvalidDayOfWeek // Assuming ErrInvalidDayOfWeek is a predefined error
	}
	if endTime.Before(startTime) {
		return nil, ErrInvalidTimeRange // Assuming ErrInvalidTimeRange is a predefined error
	}
	if breakEndTime.Before(breakStartTime) {
		return nil, ErrInvalidBreakTime // Assuming ErrInvalidBreakTime is a predefined error
	}

	return &EmployeeSchedule{
		DayOfWeek:      dayOfWeek,
		StartTime:      startTime,
		EndTime:        endTime,
		BreakStartTime: breakStartTime,
		BreakEndTime:   breakEndTime,
	}, nil
}

type EmployeeModel struct {
	DB *sql.DB
}

func (m EmployeeModel) Insert(employee *Employee) error {
	query := `Insert into employee (inst_id, user_id, description, name, photo_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	args := []any{employee.InstId, employee.UserId, employee.Description, employee.Name, employee.PhotoUrl}
	var employeeId int64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = tx.QueryRowContext(ctx, query, args...).Scan(&employeeId)
	if err != nil {
		tx.Rollback()
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.Code == "23503" {
				return ErrInvalidInstId
			}
		}
		return err
	}
	employee.ID = employeeId
	for _, schedule := range employee.Schedule {
		hs, ms, ss := schedule.StartTime.Clock()
		he, me, se := schedule.EndTime.Clock()
		hbs, mbs, sbs := schedule.BreakStartTime.Clock()
		hbe, mbe, sbe := schedule.BreakEndTime.Clock()
		query = `
		INSERT INTO employee_schedule (employee_id, day_of_week, start_time, end_time, break_start_time, break_end_time) VALUES ($1, $2, $3, $4, $5, $6)`
		args = []any{employeeId, schedule.DayOfWeek, fmt.Sprintf("%d:%d:%d", hs, ms, ss), fmt.Sprintf("%d:%d:%d", he, me, se), fmt.Sprintf("%d:%d:%d", hbs, mbs, sbs), fmt.Sprintf("%d:%d:%d", hbe, mbe, sbe)}
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, service := range employee.Services {
		query = `
		INSERT INTO employee_service (employee_id, service_id) VALUES ($1, $2)`
		args = []any{employeeId, service.ServiceId}
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			tx.Rollback()
			// foreign key error handle
			if pgerr, ok := err.(pgx.PgError); ok {
				if pgerr.Code == "23503" {
					return ErrInvalidServices
				}
			}
			return err
		}
	}
	return tx.Commit()
}

func (m EmployeeModel) GetAllForInst(instId int64) ([]*Employee, error) {
	query := `
	SELECT id, created_at, inst_id, user_id, description, name, photo_url
	FROM employee
	WHERE inst_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	employeeRows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer employeeRows.Close()
	var employees []*Employee
	for employeeRows.Next() {
		var employee Employee
		err := employeeRows.Scan(
			&employee.ID,
			&employee.CreatedAt,
			&employee.InstId,
			&employee.UserId,
			&employee.Description,
			&employee.Name,
			&employee.PhotoUrl,
		)
		if err != nil {
			return nil, err
		}
		query = `
		SELECT day_of_week, start_time, end_time, break_start_time, break_end_time
		FROM employee_schedule
		WHERE employee_id = $1`
		workHoursRows, err := m.DB.QueryContext(ctx, query, employee.ID)
		if err != nil {
			return nil, err
		}
		defer workHoursRows.Close()
		for workHoursRows.Next() {
			var schedule EmployeeSchedule
			var startTime, endTime, breakStartTime, breakEndTime string
			var dayOfWeek int
			err := workHoursRows.Scan(
				&dayOfWeek,
				&startTime,
				&endTime,
				&breakStartTime,
				&breakEndTime,
			)
			if err != nil {
				return nil, err
			}
			err = schedule.ParseToTimeStruct(dayOfWeek, startTime, endTime, breakStartTime, breakEndTime)
			if err != nil {
				return nil, err
			}
			employee.Schedule = append(employee.Schedule, &schedule)
		}
		if err = workHoursRows.Err(); err != nil {
			return nil, err
		}
		query = `
		SELECT e.service_id, s.name
		FROM employee_service e
		JOIN services s on e.service_id = s.id
		WHERE employee_id = $1`
		servicesRows, err := m.DB.QueryContext(ctx, query, employee.ID)
		if err != nil {
			return nil, err
		}
		defer servicesRows.Close()
		for servicesRows.Next() {
			var service EmployeeServices
			err := servicesRows.Scan(
				&service.ServiceId,
				&service.Name,
			)
			if err != nil {
				return nil, err
			}
			employee.Services = append(employee.Services, &service)
		}
		if err = servicesRows.Err(); err != nil {
			return nil, err
		}
		employees = append(employees, &employee)
	}
	if err = employeeRows.Err(); err != nil {
		return nil, err
	}
	return employees, nil
}

func (m EmployeeModel) GetById(id int64) (*Employee, error) {
	query := `
	SELECT id, created_at, inst_id, user_id, description, name, photo_url
	FROM employee
	WHERE id = $1`
	var employee Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&employee.ID,
		&employee.CreatedAt,
		&employee.InstId,
		&employee.UserId,
		&employee.Description,
		&employee.Name,
		&employee.PhotoUrl,
	)
	if err != nil {
		return nil, err
	}
	query = `
	SELECT day_of_week, start_time, end_time, break_start_time, break_end_time
	FROM employee_schedule
	WHERE employee_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var schedule EmployeeSchedule
		err := rows.Scan(
			&schedule.DayOfWeek,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.BreakStartTime,
			&schedule.BreakEndTime,
		)
		if err != nil {
			return nil, err
		}
		employee.Schedule = append(employee.Schedule, &schedule)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	query = `
		SELECT e.service_id, s.name
		FROM employee_service e
		JOIN service s on e.service_id = s.id
		WHERE employee_id = $1`
	rows, err = m.DB.QueryContext(ctx, query, employee.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var service EmployeeServices
		err := rows.Scan(
			&service.ServiceId,
			&service.Name,
		)
		if err != nil {
			return nil, err
		}
		employee.Services = append(employee.Services, &service)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &employee, nil
}

func (m EmployeeModel) Update(employee *Employee) error {
	query := `
	UPDATE employee
	SET description = $1,
		photo_url = $2,
		name = $3
	WHERE id = $4`
	args := []any{employee.Description, employee.PhotoUrl, employee.Name, employee.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m EmployeeModel) Delete(id int64) error {
	query := `
	DELETE FROM employee
	WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m EmployeeModel) UpdateSchedule(employee *Employee) error {
	query := `
	DELETE FROM employee_schedule
	WHERE employee_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, employee.ID)
	if err != nil {
		return err
	}
	for _, schedule := range employee.Schedule {
		query = `
		INSERT INTO employee_schedule (employee_id, day_of_week, start_time, end_time, break_start_time, break_end_time) VALUES ($1, $2, $3, $4, $5, $6)`
		args := []any{employee.ID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime, schedule.BreakStartTime, schedule.BreakEndTime}
		_, err = m.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m EmployeeModel) UpdateServices(employee *Employee) error {
	query := `
	DELETE FROM employee_service
	WHERE employee_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, employee.ID)
	if err != nil {
		return err
	}
	for _, service := range employee.Services {
		query = `
		INSERT INTO employee_service (employee_id, service_id) VALUES ($1, $2)`
		args := []any{employee.ID, service.ServiceId}
		_, err = m.DB.ExecContext(ctx, query, args...)
		if err != nil {
			if pgerr, ok := err.(pgx.PgError); ok {
				if pgerr.Code == "23503" {
					return ErrRecordNotFound
				}
			}
			return err
		}
	}
	return nil
}

func (m EmployeeModel) GetEmployeesForInstitution(instId int64) ([]*Employee, error) {
	query := `
	SELECT id, created_at, inst_id, user_id, description, name, photo_url
	FROM employee
	WHERE inst_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []*Employee
	for rows.Next() {
		var employee Employee
		err := rows.Scan(
			&employee.ID,
			&employee.CreatedAt,
			&employee.InstId,
			&employee.UserId,
			&employee.Description,
			&employee.Name,
			&employee.PhotoUrl,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, &employee)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	query = `SELECT day_of_week, start_time, end_time, break_start_time, break_end_time
    	FROM employee_schedule
		WHERE employee_id = $1`
	for _, employee := range employees {
		rows, err := m.DB.QueryContext(ctx, query, employee.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var schedule EmployeeSchedule
			err := rows.Scan(
				&schedule.DayOfWeek,
				&schedule.StartTime,
				&schedule.EndTime,
				&schedule.BreakStartTime,
				&schedule.BreakEndTime,
			)
			if err != nil {
				return nil, err
			}
			employee.Schedule = append(employee.Schedule, &schedule)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	query = `SELECT e.service_id, s.name
    	FROM employee_service e
    			JOIN services s on e.service_id = s.id
    					WHERE employee_id = $1`
	for _, employee := range employees {
		rows, err := m.DB.QueryContext(ctx, query, employee.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var service EmployeeServices
			err := rows.Scan(
				&service.ServiceId,
				&service.Name,
			)
			if err != nil {
				return nil, err
			}
			employee.Services = append(employee.Services, &service)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return employees, nil
}

func (m EmployeeModel) GetEmployeeScheduleAndService(employeeId int64, serviceId int64, selectedDay time.Time) (*TypeForEmployeeTimeSlots, error) {
	query := `
	SELECT day_of_week, start_time, end_time, break_start_time, break_end_time
	FROM employee_schedule
	WHERE employee_id = $1 AND day_of_week = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var schedule ScheduleForEmployeeTimeSlots
	err := m.DB.QueryRowContext(ctx, query, employeeId, int(selectedDay.Weekday())).Scan(
		&schedule.DayOfWeek,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.BreakStartTime,
		&schedule.BreakEndTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	
	query = `
	SELECT name, duration
	FROM services
	WHERE id = $1`
	var service ServiceForEmployeeTimeSlots
	var durationStr string
	err = m.DB.QueryRowContext(ctx, query, serviceId).Scan(&service.Name, &durationStr)
	if err != nil {
		return nil, err
	}
	durationStr = strings.Replace(durationStr, ":", "h", 1)
	durationStr = strings.Replace(durationStr, ":", "m", 1)
	service.Duration = durationStr
	return &TypeForEmployeeTimeSlots{
		Schedule: schedule,
		Service:  service,
	}, nil
}
