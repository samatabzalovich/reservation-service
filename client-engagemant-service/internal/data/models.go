package data

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidField   = errors.New("invalid field, check the fields you entered and try again")
	ErrRecordNotFound = errors.New("record not found")
)

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Rating:  &RatingModel{DB: db},
		Comment: &CommentModel{DB: db},
	}
}

type Models struct {
	Rating interface {
		Insert(rating *Rating) error
		GetAllForInst(instId int64) ([]*Rating, error)
		GetById(id int64) (*Rating, error)
		Update(rating *Rating) error
		Delete(id int64) error
		GetAllForClient(clientId int64) ([]*Rating, error)
		GetAllForEmployee(employeeId int64) ([]*Rating, error)
		GetRatingForAppointment(appointmentId int64) (*Rating, error)
		GetRatingsForEmployee(employeeId int64) ([]*Rating, error)
		GetRatingsForClient(clientId int64) ([]*Rating, error)
		GetRatingsForInstitution(instId int64) ([]*Rating, error)
		GetAverageRatingForEmployee(employeeId int64) (float64, error)
		GetAverageRatingForInstitution(instId int64) (float64, error)
		GetAverageRatingForService(serviceId int64) (float64, error)
		GetAverageRatingForClient(clientId int64) (float64, error)
		GetAverageRatingForAppointment(appointmentId int64) (float64, error)
		GetAverageRatingForEmployeeService(employeeId, serviceId int64) (float64, error)
		GetAverageRatingForClientService(clientId, serviceId int64) (float64, error)
		GetAverageRatingForClientInstitution(clientId, instId int64) (float64, error)
		GetAverageRatingForClientEmployee(clientId, employeeId int64) (float64, error)
	}
	Comment interface {
		Insert(c *Comment) error
		GetById(id int64) (*Comment, error)
		GetByInstitutionId(instId int64) ([]*Comment, error)
		GetByUserId(userId int64) ([]*Comment, error)
		Update(c *Comment) error
		Delete(id int64) error
		DeleteByInstitutionId(instId int64) error
		DeleteByUserId(userId int64) error
	}
}
