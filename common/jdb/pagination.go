package jdb

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())

	ErrInvalidOrderBy = jerror.NewValidationError("INVALID_ORDER_BY")
	ErrInvalidOrder   = jerror.NewValidationError("INVALID_ORDER")
)

type QuerySort struct {
	OrderBy string `form:"order_by,omitempty" swaggerignore:"true"` // swaggerignore because we have to define the allowed values manually
	Order   string `form:"order,omitempty" binding:"omitempty,oneof=asc desc ASC DESC"`
}

func (qs QuerySort) GetOrderBy(validator string, defaultOrderBy string) (string, error) {
	if validator != "" {
		if err := validate.Var(qs.OrderBy, validator); err != nil {
			return "", fmt.Errorf("%s: %w", err.Error(), ErrInvalidOrderBy)
		}
	}
	if err := validate.Var(qs.Order, "omitempty,oneof=asc desc ASC DESC"); err != nil {
		return "", fmt.Errorf("%s: %w", err.Error(), ErrInvalidOrder)
	}

	if qs.OrderBy == "" && defaultOrderBy != "" {
		qs.OrderBy = defaultOrderBy
	} else if qs.OrderBy == "" {
		return "", nil
	}

	order := "DESC"
	if qs.Order != "" {
		order = strings.ToUpper(qs.Order)
	}
	return qs.OrderBy + " " + order, nil
}

func (qs QuerySort) OrderByMapping(mapping map[string][]string, defaultOrderBy string) (string, error) {
	if qs.OrderBy == "" {
		if defaultOrderBy == "" {
			return "", nil
		}
		qs.OrderBy = defaultOrderBy
	}

	if mapping == nil {
		return "", fmt.Errorf("invalid order by %q: %w", qs.OrderBy, ErrInvalidOrderBy)
	}

	orderBy, ok := mapping[qs.OrderBy]
	if !ok || len(orderBy) == 0 {
		return "", fmt.Errorf("invalid order by %q: %w", qs.OrderBy, ErrInvalidOrderBy)
	}

	order := "DESC"
	if qs.Order != "" {
		order = strings.ToUpper(qs.Order)
		if order != "ASC" && order != "DESC" {
			return "", fmt.Errorf("invalid order %q: %w", qs.Order, ErrInvalidOrder)
		}
	}

	var statement strings.Builder
	for i, ob := range orderBy {
		statement.WriteString(ob)
		statement.WriteString(" ")
		statement.WriteString(order)
		statement.WriteString(" NULLS LAST")
		if i != len(orderBy)-1 {
			statement.WriteString(", ")
		}
	}

	return statement.String(), nil
}

// QueryPagination is a default struct used to handle basic pagination using query params.
type QueryPagination struct {
	QuerySort
	Offset uint64 `form:"offset,omitempty"`
	Limit  uint64 `form:"limit,omitempty"`
}

// ApplyPaginationSq add the pagination specs to the given squirrel select builder.
// The validator represents the https://pkg.go.dev/github.com/go-playground/validator/v10 used to validate the order by field.
func (qp QueryPagination) ApplyPaginationSq(builder squirrel.SelectBuilder, validator string, defaultOrderBy string) (squirrel.SelectBuilder, error) {
	if qp.Limit > 0 {
		builder = builder.Limit(qp.Limit)
	}
	if qp.Offset > 0 {
		builder = builder.Offset(qp.Offset)
	}

	orderBy, err := qp.GetOrderBy(validator, defaultOrderBy)
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}
	if orderBy == "" {
		return builder, nil
	}
	return builder.OrderBy(orderBy), nil
}

func (qp QueryPagination) AddURLValues(uri url.URL) url.URL {
	values := uri.Query()

	if qp.Offset != 0 {
		values.Add("offset", strconv.FormatUint(qp.Offset, 10))
	}
	if qp.Limit != 0 {
		values.Add("limit", strconv.FormatUint(qp.Limit, 10))
	}
	if qp.Order != "" {
		values.Add("order_by", qp.OrderBy)
	}
	if qp.Order != "" {
		values.Add("order", qp.Order)
	}

	uri.RawQuery = values.Encode()
	return uri
}

func (qp QueryPagination) IsDesc() bool {
	return strings.ToLower(qp.Order) == "desc"
}
