package data

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidToken       = errors.New("invalid auth token")
	ErrInvalidPageSize    = errors.New("invalid page size")
	ErrInvalidPage        = errors.New("invalid page number")
	ErrInvalidQueueStatus = errors.New("invalid queue status")
	ErrInvalidID          = errors.New("invalid id")
	ErrRecordNotFound     = errors.New("record not found")
	ErrNoClientInQueue    = errors.New("no client in queue")
	ErrConcurrentUpdate   = errors.New("concurrent update")
	ErrInvalidQueueInfo   = errors.New("invalid queue info")
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
		Insert(clientId, institutionID, serviceID int64) (*Queue, error)
		GetAllForInst(instId int64, pageInfo Filters) ([]*Queue, Metadata, error)
		GetById(id int64) (*Queue, error)
		Update(queue *Queue) error
		Delete(id int64) error
		DeleteAllForInst(instId int64) error
		GetForClient(clientId int64) ([]*Queue, Metadata, error)
		CallFromQueue(queue *Queue) error
		GetLastPositionedQueue(serviceId int64) (*Queue, error)
		GetLastForClient(clientId int64) (*Queue, error)
	}
}
