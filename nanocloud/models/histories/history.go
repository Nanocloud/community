package histories

type History struct {
	Id           string `json:"-"`
	UserId       string `json:"user-id"`
	ConnectionId string `json:"connection-id"`
	StartDate    string `json:"start-date"`
	EndDate      string `json:"end-date"`
}

func (h *History) GetID() string {
	return h.Id
}

func (h *History) SetID(id string) error {
	h.Id = id
	return nil
}
