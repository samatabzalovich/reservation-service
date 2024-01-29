package data

import (
	"context"
	"database/sql"
	"errors"
	"institution-service/internal/validator"
	"time"
)

type Institution struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	OwnerId     int64  `json:"owner_id"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Country     string `json:"country"`
	City        string `json:"city"`
	CategoryId  int64  `json:"category_id"`
	Version     int    `json:"-"`
}

func NewInstitution(id int64, name string, description string, website string, ownerId int64, latitude string, longitude string, country string, city string, categoryId int64, phone string, address string) (*Institution, error) {
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
	if city == "" {
		return nil, ErrInvalidCity
	}
	if categoryId < 1 {
		return nil, ErrInvalidCategoryId
	}
	if phone == "" || validator.PhoneRX.MatchString(phone) == false {
		return nil, ErrInvalidPhone
	}
	if address == "" {
		return nil, ErrInvalidAddress
	}
	return &Institution{
		ID:          id,
		Name:        name,
		Description: description,
		Website:     website,
		OwnerId:     ownerId,
		Latitude:    latitude,
		Longitude:   longitude,
		Country:     country,
		City:        city,
		CategoryId:  categoryId,
		Phone:       phone,
	}, nil
}

type InstitutionModel struct {
	DB *sql.DB
}

func(m InstitutionModel) Insert(institution * Institution)(int64, error) {
        query := `INSERT INTO institutions (name, description, website, owner_id, latitude, longitude, country, city, category_id, phone, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,$11) RETURNING id` 
		ctx,cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
		args := []interface{} {
        institution.Name,
        institution.Description,
        institution.Website,
        institution.OwnerId,
        institution.Latitude,
        institution.Longitude,
        institution.Country,
        institution.City,
        institution.CategoryId,
        institution.Phone,
        institution.Address,
    }
    var id int64 
	err := m.DB.QueryRowContext(ctx, query, args ...).Scan(&id)
	if err != nil {
		return 0, err
    }
    return id,nil
}

func (m InstitutionModel) GetById(id int64) (*Institution, error) {
	query := `SELECT id, name, description, website, owner_id, latitude, longitude, country, city, category_id, phone, address FROM institutions WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var institution *Institution
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Country, &institution.City, &institution.CategoryId, &institution.Phone, &institution.Address)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m InstitutionModel) GetAll(categoryId int64, filters Filters) ([]*Institution, Metadata, error) {
	query := `SELECT count(*) OVER(), id, name, description, website, owner_id, latitude, longitude, country, city, category_id, phone, address FROM institutions WHERE category_id = $1 ORDER BY ` + filters.Sort + ` LIMIT $2 OFFSET $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{categoryId, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	var institutions []*Institution
	totalRecords := 0
	for rows.Next() {
		var institution Institution
		err = rows.Scan(&totalRecords, &institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Country, &institution.City, &institution.CategoryId, &institution.Phone, &institution.Address)
		if err != nil {
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

func (m InstitutionModel) Update(institution *Institution) error {
	query := `UPDATE institutions SET name = $1, description = $2, website = $3, owner_id = $4, latitude = $5, longitude = $6, country = $7, city = $8, category_id = $9, phone = $10, address = $11, version = version + 1 WHERE id = $12 AND version = $13 RETURNING version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{
		institution.Name,
		institution.Description,
		institution.Website,
		institution.OwnerId,
		institution.Latitude,
		institution.Longitude,
		institution.Country,
		institution.City,
		institution.CategoryId,
		institution.Phone,
		institution.Address,
		institution.ID,
		institution.Version,
	}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&institution.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m InstitutionModel) Search(categoryId int64, searchText string, filters Filters) ([]*Institution, Metadata, error) {
	query := `SELECT count(*) OVER(), id, name, description, website, owner_id, latitude, longitude, country, city, category_id, phone, address FROM institutions WHERE (category_id = $1) AND (to_tsvector('simple', name) @@ plainto_tsquery('simple', $2) OR to_tsvector('simple', description) @@ plainto_tsquery('simple', $2)) ORDER BY ` + filters.Sort + ` LIMIT $3 OFFSET $4`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{categoryId, searchText, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	var institutions []*Institution
	totalRecords := 0
	for rows.Next() {
		var institution Institution
		err = rows.Scan(&totalRecords, &institution.ID, &institution.Name, &institution.Description, &institution.Website, &institution.OwnerId, &institution.Latitude, &institution.Longitude, &institution.Country, &institution.City, &institution.CategoryId, &institution.Phone, &institution.Address)
		if err != nil {
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
	query := `DELETE FROM institutions WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
