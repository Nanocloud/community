package histories

import (
	"github.com/Nanocloud/community/nanocloud/models/apps"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/manyminds/api2go/jsonapi"
)

type History struct {
	Id           string `json:"-"`
	UserId       string `json:"user-id"`
	ConnectionId string `json:"connection-id"`
	StartDate    string `json:"start-date"`
	EndDate      string `json:"end-date"`

	user        *users.User
	application *apps.Application
}

func (h *History) GetID() string {
	return h.Id
}

func (h *History) SetID(id string) error {
	h.Id = id
	return nil
}

func (h *History) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "user",
		},
		{
			Type: "applications",
			Name: "application",
		},
	}
}

func (h *History) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if h.user != nil {
		result = append(
			result,
			jsonapi.ReferenceID{
				ID:   h.user.Id,
				Name: "user",
				Type: "users",
			},
		)
	}

	if h.application != nil {
		result = append(
			result,
			jsonapi.ReferenceID{
				ID:   h.application.GetID(),
				Name: "application",
				Type: "applications",
			},
		)
	}

	return result
}

func (h *History) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if h.user != nil {
		result = append(result, h.user)
	}

	if h.application != nil {
		result = append(result, h.application)
	}

	return result
}
