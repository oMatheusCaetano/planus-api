package dto

type CreatePersonDTO struct {
    Name  string `json:"name" binding:"required,filled,min=2,max=255"`
}

type UpdatePersonDTO struct {
    Name  string `json:"name" binding:"omitempty,filled,min=2,max=255"`
}
