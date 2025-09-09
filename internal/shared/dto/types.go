package dto

type PaginationSortBy struct {
    Key       string `json:"key"`
    Direction string `json:"direction"`
}

type PaginationWhere struct {
    Key      string      `json:"key"`
    Operator string      `json:"operator"`
    Type     string      `json:"type"`
    Value    interface{} `json:"value"`
}

type PaginationMeta struct {
    Total       int                `json:"total"`
    PerPage     int                `json:"per_page"`
    CurrentPage int                `json:"current_page"`
    LastPage    int                `json:"last_page"`
    FirstPage   int                `json:"first_page"`
    SortBy      []PaginationSortBy `json:"sort_by"`
    Where       []PaginationWhere  `json:"where"`
}

type Paginated[T any] struct {
	Meta PaginationMeta `json:"meta"`
    Data []T            `json:"data"`
}
