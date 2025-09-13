package dto

type SortBy struct {
    Key       string `json:"key"`
    Direction string `json:"direction"`
}

type Where struct {
	Key      string  `json:"key"`
	Operator string  `json:"operator"`
	Value    any     `json:"value"`
}

type WhereLogicBlock struct {
	Operator  string              `json:"operator"` // "and" or "or"
	Condition *Where              `json:"condition,omitempty"`
	Sub       []*WhereLogicBlock  `json:"sub,omitempty"`
}

type PaginatedData[T any] struct {
    Data  []*T          `json:"data"`
    Meta  *PaginationMeta `json:"meta"`
}

type Paginate struct {
    Page    int                `json:"page"     binding:"omitempty,min=1"`
    PerPage int                `json:"per_page" binding:"omitempty,min=1"`
    SortBy  []*SortBy          `json:"sort_by"  binding:"omitempty,dive"`
    Where   []*WhereLogicBlock `json:"where"    binding:"omitempty,dive"`
}

type PaginationMeta struct {
    Page      int                    `json:"page"`
    PerPage   int                    `json:"per_page"`
    LastPage  int                    `json:"last_page"`
    Total     int                    `json:"total"`
    SortBy    []*SortBy              `json:"sort_by"`
    Where     []*WhereLogicBlock     `json:"where"`
}

type List struct {
    SortBy []*SortBy          `json:"sort_by" binding:"omitempty,dive"`
    Where  []*WhereLogicBlock `json:"where"   binding:"omitempty,dive"`
}

