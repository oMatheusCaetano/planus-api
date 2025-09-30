package dto

type CreateUserDTO struct {
	Name 	  string `json:"name"     binding:"required,min=2,max=255"`
	Username  string `json:"username" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=255"`
}

type UpdateUserDTO struct {
	Name 	  string `json:"name"     binding:"omitempty,min=2,max=255"`
	Username  string `json:"username" binding:"omitempty,email"`
}
