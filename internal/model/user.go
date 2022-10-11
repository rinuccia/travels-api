package model

// User represent user data model
type User struct {
	UserId    uint32 `json:"user_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Gender    string `json:"gender" validate:"required,eq=f|eq=m"`
}
