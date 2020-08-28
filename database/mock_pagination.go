package database

import (
	"errors"
	"math"
	"sort"
	"time"

	"learning/unit-testing/models"

	"github.com/go-openapi/strfmt"
)

// PaginatedQuery - that holds all the necessary info to make a paginated query.
type PaginatedQuery struct {
	CollectionName       string
	UserConnectionGroups []internal.UserConnectionGroupInfo
	Offset               int
	Limit                int
	OrderBy              string
	Order                string
	ResultCount          int
}

// NewPaginatedQuery -
func (pq *PaginatedQuery) SetPaginatedQuery(offset *int32, limit *int32, orderBy *string, order *string) {

	// Set the offset.
	if internal.IsZeroOfUnderlyingType(offset) {
		pq.Offset = 0
	} else {
		pq.Offset = int(*offset)
	}

	// Set the limit.
	if !internal.IsZeroOfUnderlyingType(limit) {
		pq.Limit = int(*limit)
	}

	// Set the orderBy.
	if !internal.IsZeroOfUnderlyingType(orderBy) {
		pq.OrderBy = string(*orderBy)
	}

	// Set the order.
	if !internal.IsZeroOfUnderlyingType(order) {
		pq.Order = string(*order)
	}
}

// SortPaginatedQuery -
func (pq *PaginatedQuery) SortPaginatedQuery() {
	switch pq.CollectionName {
	case "users_connections_groups":
		// Sort order by asc, desc
		switch pq.Order {
		case "asc":
			sort.SliceStable(pq.UserConnectionGroups, func(i, j int) bool {
				return pq.UserConnectionGroups[i].GroupName < pq.UserConnectionGroups[j].GroupName
			})
		case "desc":
			sort.SliceStable(pq.UserConnectionGroups, func(i, j int) bool {
				return pq.UserConnectionGroups[i].GroupName > pq.UserConnectionGroups[j].GroupName
			})
		}
	}
}

// AddTimeSetFilterToQuery - Adds time sets query parameters to a query. Checks for an invalid time range.
func (pq *PaginatedQuery) AddTimeSetFilterToQuery(fieldName string, afterTime *strfmt.DateTime, beforeTime *strfmt.DateTime) error {

	if !internal.IsZeroOfUnderlyingType(afterTime) {
		if !internal.IsZeroOfUnderlyingType(beforeTime) {
			if time.Time(*afterTime).After(time.Time(*beforeTime)) {
				return errors.New("invalid time range: after time can't come after before time")
			}

			switch pq.CollectionName {
			case "users_connections_groups":
				var groups []internal.UserConnectionGroupInfo
				for _, row := range pq.UserConnectionGroups {
					if time.Time(row.LatestInteractionTime).After(time.Time(*afterTime)) && time.Time(row.LatestInteractionTime).Before(time.Time(*beforeTime)) {
						groups = append(groups, row)
					}
				}
				pq.UserConnectionGroups = groups
			}
		}
	}

	return nil
}

// LimitPaginatedQuery -
func (pq *PaginatedQuery) LimitPaginatedQuery() {

	switch pq.CollectionName {
	case "users_connections_groups":

		pq.ResultCount = len(pq.UserConnectionGroups)

		// Filter data based on limit and offset
		var groups []internal.UserConnectionGroupInfo
		limit := pq.Limit + pq.Offset
		if limit > len(pq.UserConnectionGroups) {
			limit = len(pq.UserConnectionGroups)
		}
		for i := pq.Offset; i < limit; i++ {
			groups = append(groups, pq.UserConnectionGroups[i])
		}
		pq.UserConnectionGroups = groups
	}

}

// GetPaginatedQueryMetadata - Get the count of the query as filtered.
func (pq *PaginatedQuery) GetPaginatedQueryMetadata() *models.PaginationData {

	// Get the total count.
	limit := pq.Limit
	offset := pq.Offset
	resultCount := int32(pq.ResultCount)

	// Calculate the pages.
	pageCount := int32(math.Ceil(float64(resultCount) / float64(limit)))

	currentPage := int32(1)

	if offset == 0 {
		currentPage = 1
	} else {
		currentPage = int32(math.Ceil((float64(offset) + 1) / float64(pq.Limit)))
	}

	limit32 := int32(limit)

	// Create the pagination data.
	var paginationInfo = models.PaginationData{
		ResultCount: &resultCount,
		PageLimit:   &limit32,
		PageCount:   &pageCount,
		CurrentPage: &currentPage,
	}

	return &paginationInfo
}
