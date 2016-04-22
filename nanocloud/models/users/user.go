package users

type User struct {
	Id              string `json:"-"`
	Email           string `json:"email"`
	Password        string `json:"password,omitempty"`
	Activated       bool   `json:"activated"`
	IsAdmin         bool   `json:"is-admin"`
	FirstName       string `json:"first-name"`
	LastName        string `json:"last-name"`
	Sam             string `json:"sam"`
	WindowsPassword string `json:"windows-password"`
}

func (u *User) GetID() string {
	return u.Id
}

func (u *User) SetID(id string) error {
	u.Id = id
	return nil
}
