package users

import (
	"errors"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
)

type WindowsUser struct {
	Sam      string
	Password string
	Domain   string
}

type User struct {
	Id         string `json:"-"`
	Email      string `json:"email"`
	Password   string `json:"password,omitempty"`
	Activated  bool   `json:"activated"`
	IsAdmin    bool   `json:"is-admin"`
	FirstName  string `json:"first-name"`
	LastName   string `json:"last-name"`
	SignupDate int    `json:"signup-date,omitempty"`
}

func (u *User) GetID() string {
	return u.Id
}

func (u *User) SetID(id string) error {
	u.Id = id
	return nil
}

func (u *User) WindowsCredentials() (*WindowsUser, error) {
	res, err := db.Query(
		`SELECT
			windows_users.sam,
			windows_users.password,
			windows_users.domain
		FROM users_windows_user
		LEFT JOIN
			windows_users
			ON users_windows_user.windows_user_id = windows_users.id
		WHERE users_windows_user.user_id = $1::varchar`,
		u.Id,
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if !res.Next() {
		return nil, errors.New("No credentials found for this user")
	}

	winUser := WindowsUser{}
	res.Scan(
		&winUser.Sam,
		&winUser.Password,
		&winUser.Domain,
	)
	return &winUser, nil
}
