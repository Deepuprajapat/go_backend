package response

import (
	"github.com/VI-IM/im_backend_go/ent"
)

type Location struct {
	ID           string `json:"id"`
	LocalityName string `json:"locality_name"`
	City         string `json:"city"`
	State        string `json:"state"`
	PhoneNumber  string `json:"phone_number"`
	Country      string `json:"country"`
	Pincode      string `json:"pincode"`
	IsActive     bool   `json:"is_active"`
}

func GetLocationFromEnt(location *ent.Location) *Location {
	return &Location{
		ID:           location.ID,
		LocalityName: location.LocalityName,
		City:         location.City,
		State:        location.State,
		PhoneNumber:  location.PhoneNumber,
		Country:      location.Country,
		Pincode:      location.Pincode,
		IsActive:     location.IsActive,
	}
}
