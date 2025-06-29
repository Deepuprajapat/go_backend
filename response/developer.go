package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
)

type Developer struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name"`
	LegalName       string                       `json:"legal_name"`
	Identifier      string                       `json:"identifier"`
	EstablishedYear int                          `json:"established_year"`
	MediaContent    schema.DeveloperMediaContent `json:"media_content"`
	IsVerified      bool                         `json:"is_verified"`
}

func GetDeveloperFromEnt(developer *ent.Developer) *Developer {
	return &Developer{
		ID:              developer.ID,
		Name:            developer.Name,
		LegalName:       developer.LegalName,
		Identifier:      developer.Identifier,
		EstablishedYear: developer.EstablishedYear,
		MediaContent:    developer.MediaContent,
		IsVerified:      developer.IsVerified,
	}
}
