package data

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidToken    = errors.New("invalid auth token")
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidPage     = errors.New("invalid page number")
)

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Queue: QueueModel{DB: db},
	}
}

type Models struct {
	Queue interface {
		Insert(queue *Queue) error
		GetAllForInst(instId int64) ([]*Queue, error)
		GetById(id int64) (*Queue, error)
		Update(queue *Queue) error
		Delete(id int64) error
		DeleteAllForInst(instId int64) error
		GetForClient(clientId int64) ([]*Queue, Metadata, error)
		CallFromQueue(serviceId int64) error
	}
}
