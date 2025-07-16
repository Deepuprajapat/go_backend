package application

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) GetProjectByID(id string) (*response.Project, *imhttp.CustomError) {

	isDeleted, err := c.repo.IsProjectDeleted(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if project is deleted")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check if project is deleted", err.Error())
	}
	if isDeleted {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Project not found or deleted", "Project not found or deleted")
	}

	project, err := c.repo.GetProjectByID(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project", err.Error())
	}
	return response.GetProjectFromEnt(project), nil
}

func (c *application) AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, *imhttp.CustomError) {

	var project domain.Project

	exist, err := c.repo.ExistDeveloperByID(input.DeveloperID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add project", err.Error())
	}
	if !exist {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Developer not found", "Developer not found")
	}

	project.ProjectID = fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(time.Now().Unix(), 10))))[:16]
	project.ProjectName = input.ProjectName
	project.ProjectType = input.ProjectType
	project.Slug = input.Slug
	project.DeveloperID = input.DeveloperID
	project.Locality = input.Locality
	project.ProjectCity = input.ProjectCity

	projectID, err := c.repo.AddProject(project)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add project", err.Error())
	}

	return &response.AddProjectResponse{
		ProjectID: projectID,
	}, nil
}

func (c *application) UpdateProject(input request.UpdateProjectRequest) (*response.Project, *imhttp.CustomError) {

	var project domain.Project
	isDeleted, err := c.repo.IsProjectDeleted(input.ProjectID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if project is deleted")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check if project is deleted", err.Error())
	}
	if isDeleted {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Project not found or deleted", "Project not found or deleted")
	}

	project.ProjectID = input.ProjectID
	project.ProjectName = input.ProjectName
	project.Status = input.Status
	project.PriceUnit = input.PriceUnit
	project.TimelineInfo = input.TimelineInfo
	project.MetaInfo = input.MetaInfo
	project.WebCards = input.WebCards
	project.LocationInfo = input.LocationInfo
	project.IsFeatured = input.IsFeatured
	project.IsPremium = input.IsPremium
	project.IsPriority = input.IsPriority
	project.IsDeleted = input.IsDeleted
	project.Description = input.Description

	updatedProject, err := c.repo.UpdateProject(project)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update project", err.Error())
	}

	return response.GetProjectFromEnt(updatedProject), nil
}

func (c *application) DeleteProject(id string) *imhttp.CustomError {

	isDeleted, err := c.repo.IsProjectDeleted(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if project is deleted")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check if project is deleted", err.Error())
	}
	if isDeleted {
		return imhttp.NewCustomErr(http.StatusBadRequest, "Project is already deleted", "Project is already deleted")
	}

	if err := c.repo.DeleteProject(id, false); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to delete project")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete project", err.Error())
	}

	return nil
}

func (c *application) ListProjects(request *request.GetAllAPIRequest) ([]*response.ProjectListResponse, *imhttp.CustomError) {

	if request.Filters == nil {
		request.Filters = make(map[string]interface{})
	}

	projects, err := c.repo.GetAllProjects(request.Filters)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to list projects")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to list projects", err.Error())
	}

	var projectResponses []*response.ProjectListResponse
	// If name filter is present, return full project details
	if _, hasNameFilter := request.Filters["name"]; hasNameFilter {
		for _, project := range projects {
			fullProject := response.GetProjectFromEnt(project)
			projectResponses = append(projectResponses, &response.ProjectListResponse{
				ProjectID:     fullProject.ProjectID,
				ProjectName:   fullProject.ProjectName,
				ShortAddress:  fullProject.LocationInfo.ShortAddress,
				IsPremium:     fullProject.IsPremium,
				Images:        fullProject.WebCards.Images,
				Configuration: fullProject.WebCards.Details.Configuration.Value,
				MinPrice:      fullProject.MinPrice,
				Sizes:         fullProject.WebCards.Details.Sizes.Value,
				Canonical:     project.Slug,
				// Add full project details
				FullDetails: fullProject,
			})
		}
	} else {
		for _, project := range projects {
			projectResponses = append(projectResponses, response.GetProjectListResponse(project))
		}
	}

	return projectResponses, nil
}

func (c *application) CompareProjects(projectIDs []string) (*response.ProjectComparisonResponse, *imhttp.CustomError) {
	if len(projectIDs) < 2 {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "At least 2 projects are required for comparison", "At least 2 projects are required for comparison")
	}

	var comparisonProjects []*response.ProjectComparison

	for _, projectID := range projectIDs {
		project, err := c.repo.GetProjectByID(projectID)
		if err != nil {
			logger.Get().Error().Err(err).Msg("Failed to get project")
			return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project", err.Error())
		}

		if project.IsDeleted {
			return nil, imhttp.NewCustomErr(http.StatusNotFound, "Project not found or deleted", "Project not found or deleted")
		}

		var developerName string
		if project.Edges.Developer != nil {
			developerName = project.Edges.Developer.Name
		}

		comparisonProject := &response.ProjectComparison{
			ProjectID:     project.ID,
			ProjectName:   project.Name,
			Description:   project.Description,
			Status:        project.Status,
			MinPrice:      project.MinPrice,
			MaxPrice:      project.MaxPrice,
			TimelineInfo:  project.TimelineInfo,
			LocationInfo:  project.LocationInfo,
			IsFeatured:    project.IsFeatured,
			IsPremium:     project.IsPremium,
			IsPriority:    project.IsPriority,
			WebCards:      project.WebCards,
			DeveloperName: developerName,
		}

		comparisonProjects = append(comparisonProjects, comparisonProject)
	}

	return &response.ProjectComparisonResponse{
		Projects: comparisonProjects,
	}, nil
}

func (c *application) GetProjectByURL(url string) (*ent.Project, *imhttp.CustomError) {
	project, err := c.repo.GetProjectByURL(url)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get project by URL")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project", err.Error())
	}

	return project, nil
}
