package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

type Project struct {
	ProjectID    string                 `json:"project_id"`
	ProjectName  string                 `json:"project_name"`
	Description  string                 `json:"description"`
	Status       enums.ProjectStatus    `json:"status"`
	MinPrice     string                    `json:"min_price"`
	MaxPrice     string                    `json:"max_price"`
	PriceUnit    string                 `json:"price_unit"`
	TimelineInfo schema.TimelineInfo    `json:"timeline_info"`
	MetaInfo     schema.SEOMeta         `json:"meta_info"`
	WebCards     schema.ProjectWebCards `json:"web_cards"`
	LocationInfo schema.LocationInfo    `json:"location_info"`
	IsFeatured   bool                   `json:"is_featured"`
	IsPremium    bool                   `json:"is_premium"`
	IsPriority   bool                   `json:"is_priority"`
}

type AddProjectResponse struct {
	ProjectID string `json:"project_id"`
}

type UpdateProjectResponse struct {
	ProjectID string `json:"project_id"`
}

type ProjectListResponse struct {
	ProjectID     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	ShortAddress  string   `json:"short_address"`
	Canonical     string   `json:"canonical"`
	Images        []string `json:"images"`
	Configuration string   `json:"configuration"`
	MinPrice      string      `json:"min_price"`
	Sizes         string   `json:"sizes"`
	IsPremium     bool     `json:"is_premium"`
	VideoURL      string   `json:"video_url"`
	FullDetails   *Project `json:"full_details,omitempty"`
}

func GetProjectFromEnt(project *ent.Project) *Project {

	return &Project{
		ProjectID:   project.ID,
		ProjectName: project.Name,
		Description: project.Description,
		Status:      project.Status,
		MinPrice:    project.MinPrice,
		MaxPrice:    project.MaxPrice,
		TimelineInfo: schema.TimelineInfo{
			ProjectLaunchDate:     project.TimelineInfo.ProjectLaunchDate,
			ProjectPossessionDate: project.TimelineInfo.ProjectPossessionDate,
		},
		MetaInfo: project.MetaInfo,
		WebCards: project.WebCards,
		LocationInfo: schema.LocationInfo{
			ShortAddress:  project.LocationInfo.ShortAddress,
			Longitude:     project.LocationInfo.Longitude,
			Latitude:      project.LocationInfo.Latitude,
			GoogleMapLink: project.LocationInfo.GoogleMapLink,
		},
		IsFeatured: project.IsFeatured,
		IsPremium:  project.IsPremium,
		IsPriority: project.IsPriority,
	}
}

func GetProjectListResponse(project *ent.Project) *ProjectListResponse {

	return &ProjectListResponse{
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		ShortAddress:  project.LocationInfo.ShortAddress,
		IsPremium:     project.IsPremium,
		Images:        project.WebCards.Images,
		Configuration: project.WebCards.Details.Configuration.Value,
		Sizes:         project.WebCards.Details.Sizes.Value,
		VideoURL:      project.WebCards.VideoPresentation.URL,
		MinPrice:      project.MinPrice,
		Canonical:     project.MetaInfo.Canonical,
	}
}
