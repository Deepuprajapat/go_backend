package domain

import (
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

type Project struct {
	ProjectID     string
	ProjectName   string
	ProjectURL    string
	Slug          string
	ProjectType   string
	Locality      string
	ProjectCity   string
	DeveloperID   string
	Description   string
	Status        enums.ProjectStatus
	MinPrice      string
	MaxPrice      string
	PriceUnit     string
	TimelineInfo  schema.TimelineInfo
	MetaInfo      schema.SEOMeta
	WebCards      schema.ProjectWebCards
	LocationInfo  schema.LocationInfo
	IsFeatured    bool
	IsPremium     bool
	IsPriority    bool
	IsDeleted     bool
	SearchContext []string
}
