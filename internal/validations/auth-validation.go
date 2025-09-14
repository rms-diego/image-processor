package validations

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserFound struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
