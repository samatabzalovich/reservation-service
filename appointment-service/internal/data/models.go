package data

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidAppointmentTime = errors.New("invalid appointment time")
	ErrInvalidClientID        = errors.New("invalid client id")
	ErrInvalidInstID          = errors.New("invalid inst id")
	ErrInvalidEmployeeID      = errors.New("invalid employee id")
	ErrInvalidServiceID       = errors.New("invalid service id")
	
)



var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Appointments: AppointmentModel{DB: db},
	}
}

type Models struct {
	Appointments interface {
		Insert(appointment *Appointment) error
		GetAllForInst(instId int64) ([]*Appointment, error)
		GetById(id int64) (*Appointment, error)
		Update(appointment *Appointment) error
		Delete(id int64) error
		GetAllForClient(clientId int64) ([]*Appointment, error)
		GetAllForEmployee(employeeId int64) ([]*Appointment, error)
	}
}
