package request

type AddProjectRequest struct {
	ProjectName string `json:"project_name"`
	ProjectURL  string `json:"project_url"`
	ProjectType string `json:"project_type"`
	Locality    string `json:"locality"`
	ProjectCity string `json:"project_city"`
	DeveloperID string `json:"developer_id"`
}
