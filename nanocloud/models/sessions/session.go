package sessions

type Session struct {
	Id          string `json:"-"`
	SessionName string `json:"session-name"`
	Username    string `json:"username"`
	State       string `json:"state"`
}

func (h *Session) GetID() string {
	return h.Id
}

func (h *Session) SetID(id string) error {
	h.Id = id
	return nil
}
