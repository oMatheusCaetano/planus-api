package dto

type CreateUser struct {
	Name 	  string `json:"name"     binding:"required,min=2,max=255"`
	Username  string `json:"username" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=255"`
}

type UpdateUser struct {
	Name 	  string `json:"name"     binding:"omitempty,min=2,max=255"`
	Username  string `json:"username" binding:"omitempty,email"`
}
