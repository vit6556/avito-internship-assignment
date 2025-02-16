package dto

type AuthRequest struct {
	Username string `form:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `form:"password" validate:"required,min=6"`
}
