package jdb

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// SearchFilterType represents a filter type.
//
// Deprecated: [SearchCriteria] is deprecated.
type SearchFilterType string

const (
	// SearchFilterTypeOr represents a or filter type.
	SearchFilterTypeOr SearchFilterType = "||"

	// SearchFilterTypeAnd represents an and filter type.
	SearchFilterTypeAnd SearchFilterType = "&&"

	// SearchFilterTypeNotEq represents a not equal filter type.
	SearchFilterTypeNotEq SearchFilterType = "!="

	// SearchFilterTypeEq represents a equal filter type.
	SearchFilterTypeEq SearchFilterType = "=="

	// SearchFilterTypeLt represents a lower than filter type.
	SearchFilterTypeLt SearchFilterType = "<"

	// SearchFilterTypeLtOrEq represents a lower than or equal filter type.
	SearchFilterTypeLtOrEq SearchFilterType = "<="

	// SearchFilterTypeGt represents a greater than filter type.
	SearchFilterTypeGt SearchFilterType = ">"

	// SearchFilterTypeGtOrEq represents a greater than or equal filter type.
	SearchFilterTypeGtOrEq SearchFilterType = ">="

	// SearchFilterTypeLike represents a like filter type.
	SearchFilterTypeLike SearchFilterType = "like"

	// SearchFilterTypeIn represents an in filter type.
	SearchFilterTypeIn SearchFilterType = "in"

	// SearchFilterTypeNotIn represents a not in filter type.
	SearchFilterTypeNotIn SearchFilterType = "not in"
)

// SearchFilter represents a search filter.
//
// Deprecated: [SearchCriteria] is deprecated.
type SearchFilter struct {
	Type    SearchFilterType `json:"type"`
	Name    string           `json:"name"`
	Value   any              `json:"value"`
	Filters []SearchFilter   `json:"filters"`
}

// Deprecated: [SearchCriteria] is deprecated.
type Order struct {
	By            string   `json:"by"`
	ArrayPosition []string `json:"array_position"`
	Type          *string  `json:"type"`
}

// SearchCriteria represents the search criteria.
//
// Deprecated: use [QueryPagination] and specific filters instead (see: SearchClients in missions for an example).
type SearchCriteria struct {
	Filter    SearchFilter `json:"filter"`
	Limit     *uint64      `json:"limit"`
	Offset    uint64       `json:"offset"`
	OrderBy   *string      `json:"order_by"`
	OrderType *string      `json:"order_type"`
	Orders    []Order      `json:"orders"`
}

func isValidColumName(name string) bool {
	for _, c := range name {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' && c != '.' {
			return false
		}
	}

	return true
}

// BuildFilters builds squirrel filters from SearchFilter.
//
// Deprecated: [SearchCriteria] is deprecated.
func BuildFilters(filter SearchFilter) sq.Sqlizer {
	if !isValidColumName(filter.Name) {
		return nil
	}

	if filter.Type == SearchFilterTypeOr {
		group := sq.Or{}
		for _, elem := range filter.Filters {
			subFilter := BuildFilters(elem)
			if subFilter != nil {
				group = append(group, subFilter)
			}
		}
		return group
	}

	if filter.Type == SearchFilterTypeAnd {
		group := sq.And{}
		for _, elem := range filter.Filters {
			subFilter := BuildFilters(elem)
			if subFilter != nil {
				group = append(group, subFilter)
			}
		}
		return group
	}

	switch filter.Type {
	case SearchFilterTypeNotEq:
		return sq.NotEq{filter.Name: filter.Value}
	case SearchFilterTypeEq:
		return sq.Eq{filter.Name: filter.Value}
	case SearchFilterTypeLt:
		return sq.Lt{filter.Name: filter.Value}
	case SearchFilterTypeLtOrEq:
		return sq.LtOrEq{filter.Name: filter.Value}
	case SearchFilterTypeGt:
		return sq.Gt{filter.Name: filter.Value}
	case SearchFilterTypeGtOrEq:
		return sq.GtOrEq{filter.Name: filter.Value}
	case SearchFilterTypeLike:
		return sq.ILike{filter.Name: filter.Value}
	case SearchFilterTypeIn:
		return sq.Eq{filter.Name: filter.Value}
	case SearchFilterTypeNotIn:
		return sq.NotEq{filter.Name: filter.Value}
	}

	return nil
}

// ApplySearchCriteria applies search criteria to a query builder.
//
// Deprecated: [SearchCriteria] is deprecated.
func ApplySearchCriteria(builder sq.SelectBuilder, criteria SearchCriteria) sq.SelectBuilder {
	builder = builder.Offset(criteria.Offset)
	if criteria.Limit != nil {
		builder = builder.Limit(*criteria.Limit)
	}

	filters := BuildFilters(criteria.Filter)
	if filters != nil {
		builder = builder.Where(filters)
	}

	if criteria.OrderBy != nil && isValidColumName(*criteria.OrderBy) {
		orderType := "ASC"
		if criteria.OrderType != nil &&
			(*criteria.OrderType == "DESC" || *criteria.OrderType == "desc" || *criteria.OrderType == "asc") {
			orderType = *criteria.OrderType
		}
		builder = builder.OrderBy(fmt.Sprintf("%s %s", *criteria.OrderBy, orderType))
	}

	for _, order := range criteria.Orders {
		if isValidColumName(order.By) {
			orderType := "ASC"
			if order.Type != nil &&
				(*order.Type == "DESC" || *order.Type == "desc" || *order.Type == "asc") {
				orderType = *order.Type
			}

			if len(order.ArrayPosition) > 0 {
				for _, ap := range order.ArrayPosition {
					if !isValidColumName(ap) {
						return builder
					}
				}

				array := fmt.Sprintf("'%s'", strings.Join(order.ArrayPosition, "', '"))
				builder = builder.OrderBy(fmt.Sprintf("array_position(ARRAY[%s]::varchar[], %s) %s", array, order.By, orderType))
			} else {
				builder = builder.OrderBy(fmt.Sprintf("%s %s", order.By, orderType))
			}
		}
	}

	return builder
}
