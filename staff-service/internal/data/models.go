package data

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidInstId      = errors.New("invalid inst_id")
	ErrInvalidUserId      = errors.New("invalid user_id")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidSchedule    = errors.New("invalid schedule")
	ErrInvalidServices    = errors.New("invalid services")
	ErrInvalidDayOfWeek   = errors.New("invalid day of week")
	ErrInvalidTimeRange   = errors.New("invalid time range")
	ErrInvalidBreakTime   = errors.New("invalid break time")
	ErrInvalidServiceId   = errors.New("invalid service id")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrInvalidDuration    = errors.New("invalid duration")
	ErrInvalidPhotoUrl    = errors.New("invalid photo url")
	ErrRecordNotFound    = errors.New("record not found")
)

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Employees: EmployeeModel{DB: db},
		Service:   ServiceModel{DB: db},
	}
}

type Models struct {
	Employees interface {
		Insert(employee *Employee) error
		GetAllForInst(instId int64) ([]*Employee, error)
		GetById(id int64) (*Employee, error)
		Update(employee *Employee) error
		UpdateServices(employee *Employee) error 
		UpdateSchedule(employee *Employee) error 
		Delete(id int64) error
	}
	Service interface {
		Insert(service *Service) error
		GetAllForInst(instId int64) ([]*Service, error)
		GetById(id int64) (*Service, error)
		Update(service *Service) error
		Delete(id int64) error
	}
}
