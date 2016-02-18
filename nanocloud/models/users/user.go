package users

type User struct {
	Id              string `json:"-"`
	Email           string `json:"email"`
	Activated       bool   `json:"activated"`
	IsAdmin         bool   `json:"is_admin"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Sam             string `json:"sam"`
	WindowsPassword string `json:"windows_password"`
}
