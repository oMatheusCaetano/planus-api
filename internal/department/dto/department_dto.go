package dto

type CreatePerson struct {
    Name  string `json:"name" binding:"required,filled,min=2,max=255"`
}

type UpdatePerson struct {
    Name  string `json:"name" binding:"omitempty,filled,min=2,max=255"`
}
