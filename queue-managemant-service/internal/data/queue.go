package data

import (
	"database/sql"
	"time"
)

type Queue struct {
	ID          int64     `json:"id"`
	ClientID    int       `json:"clientId"`
	InstId      int       `json:"instId"`
	EmployeeID  int       `json:"employeeId"`
	ServiceID   int       `json:"serviceId"`
	QueueStatus string    `json:"queueStatus"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type QueueModel struct {
	DB *sql.DB
}

func (q QueueModel) Insert(queue *Queue) error {
	return nil
}

func (q QueueModel) GetAllForInst(instId int64) ([]*Queue, error) {
	return nil, nil
}

func (q QueueModel) GetById(id int64) (*Queue, error) {
	return nil, nil
}

func (q QueueModel) Update(queue *Queue) error {
	return nil
}

func (q QueueModel) Delete(id int64) error {
	return nil
}

func (q QueueModel) DeleteAllForInst(instId int64) error {
	return nil
}

func (q QueueModel) GetForClient(clientId int64) ([]*Queue, Metadata, error) {
	return nil, Metadata{}, nil
}

func (q QueueModel) CallFromQueue(serviceId int64) error { return nil }
