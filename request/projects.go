package request

type AddProjectRequest struct {
	ProjectName string `json:"project_name" validate:"required"`
	ProjectURL  string `json:"project_url" validate:"required"`
	ProjectType string `json:"project_type" validate:"required"`
	Locality    string `json:"locality" validate:"required"`
	ProjectCity string `json:"project_city" validate:"required"`
	DeveloperID string `json:"developer_id" validate:"required"`
}

type UpdateProjectRequest struct {
}
