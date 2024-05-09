package data

import (
	"database/sql"
	"errors"
	"time"
)

type Queue struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"clientId"`
	InstId      int64     `json:"instId"`
	EmployeeID  *int64    `json:"employeeId"`
	ServiceID   int64     `json:"serviceId"`
	Position    int       `json:"position"`
	QueueStatus string    `json:"queueStatus"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Version     int       `json:"version"`
}

func NewQueue(id, clientId, instId, employeeId, serviceId int64, position int, queueStatus string, createdAt time.Time, updatedAt time.Time) (*Queue, error) {
	queueStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	if !contains(queueStatuses, queueStatus) {
		return nil, ErrInvalidQueueStatus
	}
	if id < 1 || clientId < 1 || instId < 1 || serviceId < 1 || position < 1 {
		return nil, ErrInvalidID
	}

	return &Queue{
		ID:          id,
		ClientID:    clientId,
		InstId:      instId,
		EmployeeID:  &employeeId,
		ServiceID:   serviceId,
		Position:    position,
		QueueStatus: queueStatus,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func contains(statuses []string, status string) bool {
	for _, s := range statuses {
		if s == status {
			return true
		}
	}
	return false
}

type QueueModel struct {
	DB *sql.DB
}

func (q QueueModel) Insert(clientId, institutionID, serviceID int64) (*Queue, error) {
	var lastPosition int
	err := db.QueryRow("SELECT COALESCE(MAX(position), 0) FROM queue WHERE service_id = $1", serviceID).Scan(&lastPosition)
	if err != nil {
		return nil, err
	}

	newPosition := lastPosition + 1
	var id int64
	err = db.QueryRow("INSERT INTO queue (institution_id, service_id, position, status, user_id) VALUES ($1, $2, $3, $4, $5) returning id",
		institutionID, serviceID, newPosition, "pending", clientId).Scan(&id)
	if err != nil {
		return nil, err
	}

	return NewQueue(id, clientId, institutionID, 0, serviceID, newPosition, "pending", time.Now(), time.Now())
}

func (q QueueModel) GetAllForInst(instId int64, pageInfo Filters) ([]*Queue, Metadata, error) {
	stmt := `SELECT id, user_id, institution_id, employee_id, service_id, position, status, created_at, updated_at FROM queue WHERE institution_id = $1 LIMIT $2 OFFSET $3 `
	rows, err := db.Query(stmt, instId, pageInfo.limit(), pageInfo.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	var queues []*Queue
	var totalRecords int
	for rows.Next() {
		totalRecords++
		queue := &Queue{}
		err := rows.Scan(&queue.ID, &queue.ClientID, &queue.InstId, &queue.EmployeeID, &queue.ServiceID, &queue.Position, &queue.QueueStatus, &queue.CreatedAt, &queue.UpdatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}
		queues = append(queues, queue)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return queues, calculateMetadata(totalRecords, pageInfo.Page, pageInfo.PageSize), nil
}

func (q QueueModel) GetById(id int64) (*Queue, error) {
	stmt := `SELECT id, user_id, institution_id, employee_id, service_id, position, status, created_at, updated_at, version FROM queue WHERE id = $1`
	queue := &Queue{}
	err := db.QueryRow(stmt, id).Scan(&queue.ID, &queue.ClientID, &queue.InstId, &queue.EmployeeID, &queue.ServiceID, &queue.Position, &queue.QueueStatus, &queue.CreatedAt, &queue.UpdatedAt, &queue.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return queue, nil
}

func (q QueueModel) Update(queue *Queue) error {
	stmt := `UPDATE queue SET status = $1, employee_id = $2, updated_at = now(), version = version + 1 WHERE id = $3 AND version = $4`
	_, err := db.Exec(stmt, queue.QueueStatus, queue.EmployeeID, queue.ID, queue.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrConcurrentUpdate
		}
		return err
	}
	return nil
}

func (q QueueModel) Delete(id int64) error {
	stmt := `DELETE FROM queue WHERE id = $1`
	_, err := db.Exec(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (q QueueModel) DeleteAllForInst(instId int64) error {
	stmt := `DELETE FROM queue WHERE institution_id = $1`
	_, err := db.Exec(stmt, instId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (q QueueModel) GetForClient(clientId int64) ([]*Queue, Metadata, error) {
	stmt := `SELECT id, user_id, institution_id, employee_id, service_id, position, status, created_at, updated_at FROM queue WHERE user_id = $1`
	rows, err := db.Query(stmt, clientId)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	var queues []*Queue
	var totalRecords int
	for rows.Next() {
		totalRecords++
		queue := &Queue{}
		err := rows.Scan(&queue.ID, &queue.ClientID, &queue.InstId, &queue.EmployeeID, &queue.ServiceID, &queue.Position, &queue.QueueStatus, &queue.CreatedAt, &queue.UpdatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}
		queues = append(queues, queue)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return queues, calculateMetadata(totalRecords, 1, totalRecords), nil
}

func (q QueueModel) GetLastForClient(clientId int64) (*Queue, error) {
	stmt := `SELECT id, user_id, institution_id, employee_id, service_id, position, status, created_at, updated_at FROM queue WHERE user_id = $1 AND  status != 'completed' AND status != 'cancelled' ORDER BY position DESC LIMIT 1`
	queue := &Queue{}
	err := db.QueryRow(stmt, clientId).Scan(&queue.ID, &queue.ClientID, &queue.InstId, &queue.EmployeeID, &queue.ServiceID, &queue.Position, &queue.QueueStatus, &queue.CreatedAt, &queue.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return queue, nil
}

func (q QueueModel) GetLastPositionedQueue(serviceId int64) (*Queue, error) {
	stmt := `SELECT id, user_id, institution_id, employee_id, service_id, position, status, created_at, updated_at FROM queue WHERE service_id = $1 ORDER BY position DESC LIMIT 1`
	queue := &Queue{}
	err := db.QueryRow(stmt, serviceId).Scan(&queue.ID, &queue.ClientID, &queue.InstId, &queue.EmployeeID, &queue.ServiceID, &queue.Position, &queue.QueueStatus, &queue.CreatedAt, &queue.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return queue, nil
}

func (q QueueModel) CallFromQueue(queue *Queue) error {
	stmt := `UPDATE queue SET status = 'called' , employee_id = $1, updated_at = now(), version = version + 1 WHERE id = $2 AND version = $3`
	_, err := db.Exec(stmt, queue.EmployeeID, queue.ID, queue.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}
