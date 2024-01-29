package data

import "strings"
import "math"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func NewFilters(page int, pageSize int, sort string, sortSafelist []string) (Filters, error) {
	if page < 1 {
		return Filters{}, ErrInvalidPage
	}
	if pageSize < 1 {
		return Filters{}, ErrInvalidPageSize
	}
	if sort == "" {
		return Filters{}, ErrInvalidSort
	}
	if PermittedValue(sort, sortSafelist...) == false {
		return Filters{}, ErrInvalidSort
	}
	return Filters{
		Page:         page,
		PageSize:     pageSize,
		Sort:         sort,
		SortSafelist: sortSafelist,
	}, nil
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}