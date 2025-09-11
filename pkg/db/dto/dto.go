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
