package data

import "math"

type Filters struct {
	Page     int
	PageSize int
}

func NewFilters(page int, pageSize int) (Filters, error) {
	if page < 1 {
		return Filters{}, ErrInvalidPage
	}
	if pageSize < 1 {
		return Filters{}, ErrInvalidPageSize
	}
	return Filters{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
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
