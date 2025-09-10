package dto

import "github.com/omatheuscaetano/planus-api/internal/db"

type ListingProps struct {
    SortBy []db.SortBy
}

type PaginationProps struct {
    PerPage int
    Page   int
    SortBy []db.SortBy
}

type PaginationMeta struct {
    Total       int        `json:"total"`
    PerPage     int        `json:"per_page"`
    Page        int        `json:"page"`
    LastPage    int        `json:"last_page"`
    FirstPage   int        `json:"first_page"`
    SortBy      []db.SortBy `json:"sort_by"`
    Where       []db.Where  `json:"where"`
}

type Paginated[T any] struct {
	Meta PaginationMeta `json:"meta"`
    Data []T            `json:"data"`
}
    
