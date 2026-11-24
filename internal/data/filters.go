package data

import (
	"math"
	"strings"

	"github.com/Numeez/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type Metadata struct {
	Page         int `json:"page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records ,omitempty"`
}

func ValidateFilters(validator *validator.Validator, f Filters) {
	validator.Check(f.Page > 0, "page", "must be greater than zero")
	validator.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	validator.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	validator.Check(f.PageSize < 100, "page_size", "must be a maximum of 100")
	validator.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

func (f Filters) sortColumn() string {
	for _, val := range f.SortSafeList {
		if f.Sort == val {
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
func calculateMetdata(totalRecords, pageSize, page int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		Page:         page,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
