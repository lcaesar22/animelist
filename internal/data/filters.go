package data

import (
	"strings"

	"github.com/konaha-gakure/otaku/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	CurrentPage int `json:"current_page,omitzero"`
	PageSize    int `json:"page_size,omitzero"`
	FirstPage   int `json:"first_page,omitzero"`
	LastPage    int `json:"last_page,omitzero"`
	TotalPage   int `json:"total_page,omitzero"`
}

// the calculateMetadata() function calculates the appropriate pagination metadata values
func calculateMetadata(totalPage, page, pageSize int) Metadata {
	if totalPage == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage: page,
		PageSize:    pageSize,
		FirstPage:   1,
		LastPage:    (totalPage + pageSize - 1) / pageSize,
		TotalPage:   totalPage,
	}
}

// validate query
func ValidateFilters(v *validator.Validator, f Filters) {
	// check that the page and page_size contain the same value
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 1000000, "page", "must be greater than or equal to 1000000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be greater than 100")

	// check that the sort parameter matches a value in the safelist
	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

// Check that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
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
