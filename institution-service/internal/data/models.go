package data

import (
	"database/sql"
	"errors"
)

var (
	// ErrInvalidName is used when the name is not valid
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidDescription is used when the description is not valid
	ErrInvalidDescription = errors.New("invalid description")
	// ErrInvalidWebsite is used when the website is not valid
	ErrInvalidWebsite = errors.New("invalid website")
	// ErrInvalidOwnerId is used when the owner id is not valid
	ErrInvalidOwnerId = errors.New("invalid owner id")
	// ErrInvalidLatitude is used when the latitude is not valid
	ErrInvalidLatitude = errors.New("invalid latitude")
	// ErrInvalidLongitude is used when the longitude is not valid
	ErrInvalidLongitude = errors.New("invalid longitude")
	// ErrInvalidCountry is used when the country is not valid
	ErrInvalidCountry = errors.New("invalid country")
	// ErrInvalidCity is used when the city is not valid
	ErrInvalidCity = errors.New("invalid city")
	// ErrInvalidCategoryId is used when the category id is not valid
	ErrInvalidCategoryId = errors.New("invalid category id")
	// ErrInvalidSort is used when the sort is not valid
	ErrInvalidSort = errors.New("invalid sort")
	//ErrInvalidPageSize is used when the page size is not valid
	ErrInvalidPageSize = errors.New("invalid page size")
	//ErrInvalidPageNumber is used when the page number is not valid
	ErrInvalidPage = errors.New("invalid page number")
	//ErrInvalidPhone is used when the phone is not valid
	ErrInvalidPhone = errors.New("invalid phone")
	//ErrInvalidAddress is used when the address is not valid
	ErrInvalidAddress = errors.New("invalid address")
	//ErrEditConflict is used when the update conflicts with another update
	ErrEditConflict = errors.New("edit conflict")
	//ErrInvalidWorkingHours is used when the working hours are not valid
	ErrInvalidWorkingHours = errors.New("invalid working hours")
	//ErrInvalidDay is used when the day is not valid
	ErrInvalidDay = errors.New("invalid day")
	//ErrInvalidOpen is used when the open time is not valid
	ErrInvalidOpen = errors.New("invalid open time")
	//ErrInvalidClose is used when the close time is not valid
	ErrInvalidClose = errors.New("invalid close time")
)

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Categories:   CategoryModel{DB: db},
		Institutions: InstitutionModel{DB: db},
	}
}

type Models struct {
	Categories interface {
		Insert(category *Category) (int64, error)
		GetById(id int64) (*Category, error)
		GetAll() ([]*Category, error)
		Update(category *Category) error
		Delete(id int64) error
	}
	Institutions interface {
		Insert(institution *Institution) (int64, error)
		GetVersionByIdForOwner(ownerId, id int64) (int, error)
		GetById(id int64) (*Institution, error)
		//GetAll(categories []int64, filters Filters) ([]*Institution, Metadata, error)
		Update(institution *Institution) error
		Delete(id int64) error
		Search(categories []int64, searchText string, filters Filters) ([]*Institution, Metadata, error)
	}
}
