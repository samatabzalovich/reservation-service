package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"institution-service/internal/validator"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/lib/pq"
)

const (
	ScopeEmployeeReg = "employee_registration"
	TimeParse        = "15:04:00"
)

type Institution struct {
	ID           int64           `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Website      string          `json:"website"`
	OwnerId      int64           `json:"owner_id"`
	Latitude     string          `json:"latitude"`
	Longitude    string          `json:"longitude"`
	Address      string          `json:"address"`
	Phone        string          `json:"phone"`
	Country      string          `json:"country"`
	City         int32           `json:"city"`
	Categories   []int64         `json:"categories"`
	WorkingHours []*WorkingHours `json:"working_hours"`
	Version      int             `json:"-"`
}

type WorkingHours struct {
	DayOfWeek int       `json:"day_of_week"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}

func NewInstitution(id int64, name string, description string, website string, ownerId int64, latitude string, longitude string, country string, city int32, categories []int64, phone string, address string, workHours []*WorkingHours) (*Institution, error) {
	// validate the input
	if name == "" {
		return nil, ErrInvalidName
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if ownerId < 1 {
		return nil, ErrInvalidOwnerId
	}
	if latitude == "" {
		return nil, ErrInvalidLatitude
	}
	if longitude == "" {
		return nil, ErrInvalidLongitude
	}
	if country == "" {
		return nil, ErrInvalidCountry
	}
	if city < 1 {
		return nil, ErrInvalidCity
	}
	if len(categories) == 0 {
		return nil, ErrInvalidCategoryId
	}
	for _, id := range categories {
		if id < 1 {
			return nil, ErrInvalidCategoryId
		}
	}
	if phone == "" || !validator.PhoneRX.MatchString(phone) {
		return nil, ErrInvalidPhone
	}
	if address == "" {
		return nil, ErrInvalidAddress
	}
	if len(workHours) == 0 {
		return nil, ErrInvalidWorkingHours
	}
	return &Institution{
		ID:           id,
		Name:         name,
		Description:  description,
		Website:      website,
		OwnerId:      ownerId,
		Latitude:     latitude,
		Longitude:    longitude,
		Country:      country,
		City:         city,
		Categories:   categories,
		Phone:        phone,
		WorkingHours: workHours,
	}, nil
}

func NewWorkingHours(dayOfWeek int, openTime, closeTime string) (*WorkingHours, error) {
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return nil, ErrInvalidDay
	}
	open, err := time.Parse(TimeParse, openTime)
	if err != nil {
		return nil, ErrInvalidOpen
	}
	close, err := time.Parse(TimeParse, closeTime)
	if err != nil {
		return nil, ErrInvalidClose
	}
	return &WorkingHours{
		DayOfWeek: dayOfWeek,
		OpenTime:  open,
		CloseTime: close,
	}, nil
}

type InstitutionModel struct {
	DB *sql.DB
}

func (m InstitutionModel) Insert(institution *Institution) (int64, error) {
	query := `INSERT INTO institution (name, description, website, owner_id, latitude, longitude, country, city, phone , address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{institution.Name, institution.Description, institution.Website, institution.OwnerId, institution.Latitude, institution.Longitude, institution.Country, institution.City, institution.Phone, institution.Address}
	tx, err := m.DB.BeginTx(ctx, nil)

	if err != nil {
		return 0, err
	}
	var id int64
	err = tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, categoryId := range institution.Categories {
		query = `INSERT INTO institution_category (inst_id, category_id) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, query, id, categoryId)
		if err != nil {
			tx.Rollback()
			if pgerr, ok := err.(pgx.PgError); ok {
				if pgerr.Code == "23503" {
					return 0, ErrInvalidCategoryId
				}
			}
			return 0, err
		}
	}

	for _, workingHours := range institution.WorkingHours {
		ho, mo, so := workingHours.OpenTime.Clock()
		hc, mc, sc := workingHours.CloseTime.Clock()
		query = `INSERT INTO institution_working_hours (institution_id, day_of_week, open_time, close_time) VALUES ($1, $2, $3, $4)`
		_, err = tx.ExecContext(ctx, query, id, workingHours.DayOfWeek, fmt.Sprintf("%d:%d:%d", ho, mo, so), fmt.Sprintf("%d:%d:%d", hc, mc, sc))
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m InstitutionModel) GetVersionByIdForOwner(ownerId, id int64) (int, error) {
	query := `SELECT version FROM institution WHERE id = $1 AND owner_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var version int
	err := m.DB.QueryRowContext(ctx, query, id, ownerId).Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (m InstitutionModel) GetById(id int64) (*Institution, error) {
	query := `SELECT i.id, i.name, i.description, i.website, i.owner_id, i.latitude, i.longitude, i.address, i.phone, i.country, i.city, c.id AS category_id
	FROM institution i
	LEFT JOIN institution_category ic ON i.id = ic.inst_id
	LEFT JOIN category c ON ic.category_id = c.id
	WHERE i.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, id)
	var institution Institution
	var categoryID sql.NullInt64
	err := row.Scan(&institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Address, &institution.Phone, &institution.Country, &institution.City, &categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	if categoryID.Valid {
		institution.Categories = append(institution.Categories, categoryID.Int64)
	}

	query = `SELECT day_of_week, open_time, close_time FROM institution_working_hours WHERE institution_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var workingHours WorkingHours
		var opentime, closetime string
		err = rows.Scan(&workingHours.DayOfWeek, &opentime, &closetime)
		if err != nil {
			return nil, err
		}
		workingHours.OpenTime, err = time.Parse(TimeParse, opentime)
		if err != nil {
			return nil, ErrInvalidOpen
		}
		workingHours.CloseTime, err = time.Parse(TimeParse, closetime)
		if err != nil {
			return nil, ErrInvalidClose
		}
		institution.WorkingHours = append(institution.WorkingHours, &workingHours)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &institution, nil
}

func (m InstitutionModel) Update(institution *Institution) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	query := `UPDATE institution SET name = $1, description = $2, website = $3, owner_id = $4, latitude = $5, longitude = $6, country = $7, city = $8, phone = $9, address = $10, version = version + 1 WHERE id = $11 AND version = $12`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{institution.Name, institution.Description, institution.Website, institution.OwnerId, institution.Latitude, institution.Longitude, institution.Country, institution.City, institution.Phone, institution.Address, institution.ID, institution.Version}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	query = `DELETE FROM institution_category WHERE inst_id = $1`
	_, err = tx.ExecContext(ctx, query, institution.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `DELETE FROM institution_working_hours WHERE institution_id = $1`
	_, err = tx.ExecContext(ctx, query, institution.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, categoryId := range institution.Categories {
		query = `INSERT INTO institution_category (inst_id, category_id) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, query, institution.ID, categoryId)
		if err != nil {
			tx.Rollback()
			if pgerr, ok := err.(pgx.PgError); ok {
				if pgerr.Code == "23503" {
					return ErrInvalidCategoryId
				}
			}
			return err
		}
	}

	for _, workingHours := range institution.WorkingHours {
		ho, mo, so := workingHours.OpenTime.Clock()
		hc, mc, sc := workingHours.CloseTime.Clock()
		query = `INSERT INTO institution_working_hours (institution_id, day_of_week, open_time, close_time) VALUES ($1, $2, $3, $4)`
		_, err = tx.ExecContext(ctx, query, institution.ID, workingHours.DayOfWeek, fmt.Sprintf("%d:%d:%d", ho, mo, so), fmt.Sprintf("%d:%d:%d", hc, mc, sc))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (m InstitutionModel) Search(category []int64, searchText string, filters Filters) ([]*Institution, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), i.id, i.name, i.description, i.website,
       i.owner_id, i.latitude, i.longitude, i.address, i.phone,
       i.country, i.city, array_agg(ic.category_id) AS categories
FROM institution i
         JOIN (SELECT inst_id as id
               FROM institution_category
               WHERE (($1::int[] IS NULL OR category_id = ANY($1)))
               GROUP BY inst_id) fi ON i.id = fi.id
         LEFT JOIN institution_category ic ON i.id = ic.inst_id WHERE (to_tsvector('simple', i.name) @@ plainto_tsquery('simple', $2) 
                                          OR to_tsvector('simple', i.description) @@ plainto_tsquery('simple', $2) OR $2 = '')
GROUP BY i.id ORDER BY %s %s, id ASC LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	print(query)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, pq.Array(category), searchText, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	var institutions []*Institution
	totalRecords := 0
	for rows.Next() {
		var institution Institution
		var allCategories []int64
		err = rows.Scan(&totalRecords, &institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Address, &institution.Phone, &institution.Country, &institution.City, pq.Array(&category))
		if err != nil {
			return nil, Metadata{}, err
		}
		institution.Categories = allCategories
		query = `SELECT day_of_week, open_time, close_time FROM institution_working_hours WHERE institution_id = $1`

		workingHoursRecords, err := m.DB.QueryContext(ctx, query, institution.ID)
		if err != nil {
			return nil, Metadata{}, err
		}
		for workingHoursRecords.Next() {
			var openTime string
			var closeTime string
			var workingHours WorkingHours
			err = workingHoursRecords.Scan(&workingHours.DayOfWeek, &openTime, &closeTime)
			if err != nil {
				return nil, Metadata{}, err
			}
			workingHours.OpenTime, err = time.Parse(TimeParse, openTime)
			if err != nil {
				return nil, Metadata{}, ErrInvalidOpen
			}
			workingHours.CloseTime, err = time.Parse(TimeParse, closeTime)
			if err != nil {
				return nil, Metadata{}, ErrInvalidClose
			}
			institution.WorkingHours = append(institution.WorkingHours, &workingHours)
		}
		if err = workingHoursRecords.Err(); err != nil {
			return nil, Metadata{}, err
		}
		institutions = append(institutions, &institution)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return institutions, metadata, nil
}

func (m InstitutionModel) Delete(id int64) error {
	query := `DELETE FROM institution WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m InstitutionModel) GetForToken(tokenScope, tokenPlaintext string) (*Institution, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `SELECT i.id, i.name, i.description, i.website, i.owner_id, i.latitude, i.longitude, i.address, i.phone, i.country, i.city
	FROM institution i
	JOIN tokens ON i.owner_id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var institution Institution
	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], tokenScope, time.Now()).Scan(&institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Address, &institution.Phone, &institution.Country, &institution.City)
	log.Println("query with args: ", query, tokenHash[:], tokenScope, time.Now())
	log.Println("institution: ", institution)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &institution, nil
}

func (m InstitutionModel) GetForOwner(instId int64) ([]*Institution, Metadata, error) {
	query := `SELECT i.id, 
       i.name, i.description, 
       i.website, i.owner_id, 
       i.latitude, i.longitude, 
       i.address, i.phone, i.country,
       i.city, array_agg(ic.category_id) AS categories
FROM institution i
JOIN institution_category ic ON i.id = ic.inst_id
WHERE i.owner_id = $1
GROUP BY i.id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	var institutions []*Institution
	for rows.Next() {
		var institution Institution
		var allCategories []int64
		err = rows.Scan(&institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Address, &institution.Phone, &institution.Country, &institution.City, pq.Array(&allCategories))
		if err != nil {
			return nil, Metadata{}, err
		}
		institution.Categories = allCategories
		institutions = append(institutions, &institution)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	query = `SELECT day_of_week, open_time, close_time FROM institution_working_hours WHERE institution_id = $1`
	rows, err = m.DB.QueryContext(ctx, query, instId)
	if err != nil {
		return nil, Metadata{}, err
	}
	for _, institution := range institutions {
		workingHoursRecords, err := m.DB.QueryContext(ctx, query, institution.ID)
		if err != nil {
			return nil, Metadata{}, err
		}
		for workingHoursRecords.Next() {
			var openTime string
			var closeTime string
			var workingHours WorkingHours
			err = workingHoursRecords.Scan(&workingHours.DayOfWeek, &openTime, &closeTime)
			if err != nil {
				return nil, Metadata{}, err
			}
			workingHours.OpenTime, err = time.Parse(TimeParse, openTime)
			if err != nil {
				return nil, Metadata{}, ErrInvalidOpen
			}
			workingHours.CloseTime, err = time.Parse(TimeParse, closeTime)
			if err != nil {
				return nil, Metadata{}, ErrInvalidClose
			}
			institution.WorkingHours = append(institution.WorkingHours, &workingHours)
		}
		if err = workingHoursRecords.Err(); err != nil {
			return nil, Metadata{}, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(len(institutions), 1, len(institutions))
	return institutions, metadata, nil
}

func (m InstitutionModel) GetForEmployee(employeeId int64) (*Institution, error) {
	query := `SELECT i.id, i.name, i.description, i.website, i.owner_id, i.latitude, i.longitude, i.address, i.phone, i.country, i.city
	FROM institution i
	JOIN employee e ON i.id = e.inst_id
	WHERE e.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var institution Institution
	err := m.DB.QueryRowContext(ctx, query, employeeId).Scan(&institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Address, &institution.Phone, &institution.Country, &institution.City)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &institution, nil
}
