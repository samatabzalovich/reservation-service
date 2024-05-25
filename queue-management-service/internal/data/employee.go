package data

import (
	"database/sql"
	"errors"
)

type EmployeeInfoModel struct {
	DB *sql.DB
}

func (e EmployeeInfoModel) GetEmployeeForServiceAndUserID(serviceId, userId int64) (int64, error) {
	var instId int64
	err := e.DB.QueryRow("SELECT institution_id FROM services WHERE id = $1", serviceId).Scan(&instId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRecordNotFound
		}
		return 0, err
	}
	var employeeId int64
	err = e.DB.QueryRow("SELECT id FROM employee WHERE inst_id = $1 AND user_id = $2 LIMIT 1", instId, userId).Scan(&employeeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return 0, ErrRecordNotFound
		}
		return 0, err
	}

	return employeeId, nil
}
