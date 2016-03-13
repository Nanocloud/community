package users

type User struct {
	Id              string `json:"-"`
	Email           string `json:"email"`
	Activated       bool   `json:"activated"`
	IsAdmin         bool   `json:"isadmin"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	Sam             string `json:"sam"`
	WindowsPassword string `json:"windowspassword"`
}
