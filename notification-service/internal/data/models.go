package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound =	errors.New("record not found")
	ErrUserNotFound   = errors.New("user not found")
	ErrTokenAlreadyExists = errors.New("token already exists")
)




var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Appointments: AppointmentModel{DB: db},
		DeviceTokens: DeviceTokenModel{DB: db},
	}
}

type Models struct {
	Appointments interface {
		GetUpcomingAppointments() ([]*Appointment, error)
		MarkAsNotified(id int64) error
	}
	DeviceTokens interface {
		Insert(token DeviceToken) (int64, error)
		GetByToken(token string) (*DeviceToken, error)
		GetByUserID(userID int64) ([]*DeviceToken, error)
		Update(token DeviceToken) error
		DeleteByToken(token string) error
		Delete(id int64) error
	}
}
