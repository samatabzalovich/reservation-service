package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrInvalidField         = errors.New("invalid field, check the fields you entered and try again")
	ErrRecordNotFound       = errors.New("record not found")
	ErrRatingMustBe1        = errors.New("rating must be between 1 and 10")
	ErrInvalidAppointmentId = errors.New("invalid appointment id")
	ErrInvalidEmployeeId    = errors.New("invalid employee id")
	ErrInvalidClientId      = errors.New("invalid client id")
	ErrInvalidInstitutionId = errors.New("invalid institution id")
	ErrInvalidServiceId     = errors.New("invalid service id")
	ErrRatingAlreadyExists  = errors.New("rating already exists for this appointment")
	ErrInvalidDate          = errors.New("invalid date")
)

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Rating:    &RatingModel{DB: db},
		Comment:   &CommentModel{DB: db},
		Analytics: &AnalyticsModel{DB: db},
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
		GetAverageRatingForEmployeeService(employeeId, serviceId int64) (float64, error)
		GetAverageRatingForClientService(clientId, serviceId int64) (float64, error)
		GetAverageRatingForClientInstitution(clientId, instId int64) (float64, error)
		GetAverageRatingForClientEmployee(clientId, employeeId int64) (float64, error)
		GetServiceIdAndInstIdToRatingAppointment(appointmentId, userId int64) (*IdForRating, error)
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
	Analytics interface {
		TotalAppointmentsOfInstitutionForGivenDateRange(institutionID int64, startDate time.Time, endDate time.Time) (int, error)
		WageOfEmployeeServiceForGivenDateRange(employeeID int64, startDate time.Time, endDate time.Time) ([]*Wage, error)
		TotalRevenueOfInstitutionForGivenDateRange(institutionID int64, startDate time.Time, endDate time.Time) (Wage, error)
		TotalAppointmentsPerEmployeeForGivenDateRange(employeeID, institutionID int64, startDate, endDate time.Time) ([]*Analytics, error)
		MostPopularServicesByAppointments(serviceId, institutionId int64) ([]*Analytics, error)
		MostPopularAppointmentsBySelectedDate(institutionID int64, date string) ([]*Analytics, error)
		GetCanceledAndTotalAppointmentsForGivenDateRange(institutionID int64, startDate time.Time, endDate time.Time) ([]*Analytics, error)
		EmployeeWithHighestRating(institutionID int64) ([]*Analytics, error)
	}
}
