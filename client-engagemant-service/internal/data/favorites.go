package data

import (
	"database/sql"

	"github.com/jackc/pgx"
)

type Favorites struct {
	ID 	  int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	InstID    int64 `json:"inst_id"`
	InstName string `json:"inst_name"`
	ServiceID int64 `json:"service_id"`
	EmployeeID int64 `json:"employee_id"`
}

func NewFavorite(userID, instID, serviceID, employeeID int64) (*Favorites, error) {
	if !(userID > 0 || instID > 0 || serviceID > 0 || employeeID > 0 ){
		return nil, ErrInvalidField
	}
	f := &Favorites{
		UserID: userID,
		InstID: instID,
		ServiceID: serviceID,
		EmployeeID: employeeID,
	}
	return f, nil
}

type FavoritesModel struct {
	DB *sql.DB
}

func (m *FavoritesModel) Insert(f *Favorites) error {
	query := `INSERT INTO favorites (user_id, inst_id, service_id, employee_id)
	VALUES ($1, $2, $3, $4) RETURNING id`
	args := []interface{}{f.UserID, f.InstID, f.ServiceID, f.EmployeeID}
	err := m.DB.QueryRow(query, args...).Scan(&f.ID)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.Code == "23503" {
				return ErrInvalidField
			}
		}
		return err
	}
	return nil
}

func (m *FavoritesModel) GetById(id int64) (*Favorites, error) {
	query := `SELECT id, user_id, inst_id, service_id, employee_id FROM favorites WHERE id = $1`
	f := &Favorites{}
	err := m.DB.QueryRow(query, id).Scan(&f.ID, &f.UserID, &f.InstID, &f.ServiceID, &f.EmployeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return f, nil
}

func (m *FavoritesModel) GetByUserID(userID int64) ([]*Favorites, error) {
	query := `SELECT id, user_id, inst_id, service_id, employee_id FROM favorites WHERE user_id = $1`
	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	favorites := []*Favorites{}
	for rows.Next() {
		f := &Favorites{}
		err := rows.Scan(&f.ID, &f.UserID, &f.InstID, &f.ServiceID, &f.EmployeeID)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return favorites, nil
}

func (m *FavoritesModel) GetByInstID(instID int64) ([]*Favorites, error) {
	query := `SELECT id, user_id, inst_id, service_id, employee_id FROM favorites WHERE inst_id = $1`
	rows, err := m.DB.Query(query, instID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	favorites := []*Favorites{}
	for rows.Next() {
		f := &Favorites{}
		err := rows.Scan(&f.ID, &f.UserID, &f.InstID, &f.ServiceID, &f.EmployeeID)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return favorites, nil
}

func (m *FavoritesModel) DeleteById(id int64) error {
	query := `DELETE FROM favorites WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *FavoritesModel) DeleteByUserID(userID int64) error {
	query := `DELETE FROM favorites WHERE user_id = $1`
	_, err := m.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (m *FavoritesModel) DeleteByInstID(instID int64) error {
	query := `DELETE FROM favorites WHERE inst_id = $1`
	_, err := m.DB.Exec(query, instID)
	if err != nil {
		return err
	}
	return nil
}

func (m *FavoritesModel) DeleteByServiceID(serviceID int64) error {
	query := `DELETE FROM favorites WHERE service_id = $1`
	_, err := m.DB.Exec(query, serviceID)
	if err != nil {
		return err
	}
	return nil
}

func (m *FavoritesModel) DeleteByEmployeeID(employeeID int64) error {
	query := `DELETE FROM favorites WHERE employee_id = $1`
	_, err := m.DB.Exec(query, employeeID)
	if err != nil {
		return err
	}
	return nil
}