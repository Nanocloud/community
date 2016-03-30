package users

type User struct {
	Id              string `json:"-"`
	Email           string `json:"email"`
	Activated       bool   `json:"activated"`
	IsAdmin         bool   `json:"is-admin"`
	FirstName       string `json:"first-name"`
	LastName        string `json:"last-name"`
	Sam             string `json:"sam"`
	WindowsPassword string `json:"windows_password"`
}

func (u *User) GetID() string {
	return u.Id
}
