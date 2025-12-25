package utils

import (
	"fmt"
	"strings"

	"math"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"	
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/response"	
	"gorm.io/gorm"
)

// PaginateResult holds the results of the generic pagination operation.
type PaginateResult[T any] struct {
	Data       []T
	Pagination response.PaginationMetadata
	Error      error
}

// ApplyPagination applies sorting, searching, limit, and offset to a GORM query,
// executes the query, and calculates pagination metadata.
// T is the type of the GORM entity (e.g., entity.Cart, entity.Product).
func ApplyPagination[T any](db *gorm.DB, req request.PaginationRequest, searchAllQuery string) PaginateResult[T] {

	// --- 1. Base Query and Counting ---

	// The initial query used for counting (no offset/limit/order)
	countQuery := db.Session(&gorm.Session{})

	// The query used for fetching paginated data
	paginatedQuery := db.Session(&gorm.Session{})

	// --- 2. Apply Filtering (DRY principle applied) ---

	// 1. Specific Field Search
	if req.SearchField != "" && req.SearchValue != "" {
		whereClause := fmt.Sprintf("CAST(%s AS TEXT) ILIKE ?", req.SearchField)

		// Apply to both clones
		countQuery = countQuery.Where(whereClause, "%"+req.SearchValue+"%")
		paginatedQuery = paginatedQuery.Where(whereClause, "%"+req.SearchValue+"%")
	}

	// 2. Generic SearchAll Logic
	if req.SearchAll != "" && searchAllQuery != "" {
		search := "%" + req.SearchAll + "%"
		howManyFields := strings.Count(searchAllQuery, "?")

		if howManyFields > 0 {
			args := make([]interface{}, howManyFields)
			for i := range howManyFields {
				args[i] = search
			}

			// Apply to both clones
			countQuery = countQuery.Where(searchAllQuery, args...)
			paginatedQuery = paginatedQuery.Where(searchAllQuery, args...)
		}
	}

	// --- 3. Execute Count (Fail Fast) ---

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return PaginateResult[T]{Error: fmt.Errorf("failed to count records: %w", err)}
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	offset := (req.Page - 1) * req.PageSize

	// --- 4. Apply Pagination and Ordering ---

	// Apply Ordering
	orderBy := "id"
	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}

	sortDirection := "ASC"
	if strings.ToUpper(req.Sort) == "DESC" {
		sortDirection = "DESC"
	}
	// FAIL FAST: Use safe formatting with placeholders
	paginatedQuery = paginatedQuery.Order(fmt.Sprintf("%s %s", orderBy, sortDirection))

	// Apply Limit and Offset
	paginatedQuery = paginatedQuery.Offset(offset).Limit(req.PageSize)

	// --- 5. Execute Fetch (Fail Fast) ---

	var data []T
	if err := paginatedQuery.Find(&data).Error; err != nil {
		return PaginateResult[T]{Error: fmt.Errorf("failed to fetch records: %w", err)}
	}

	// --- 6. Build Pagination Metadata ---

	metadata := response.PaginationMetadata{
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalData: total,
		TotalPage: totalPages,
	}

	return PaginateResult[T]{
		Data:       data,
		Pagination: metadata,
	}
}
