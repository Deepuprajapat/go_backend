package handlers

import (
	"fmt"
	"net/http"
	"strings"

	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (h *Handler) GetContentTest(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// Get the request URI
	uri := r.RequestURI
	logger.Get().Info().Msgf("URL: %s", uri)

	// Build canonical URL
	canonicalURL := fmt.Sprintf("https://investmango.com%s", uri)

	// Get project by URL
	project, err := h.app.GetProjectByURL(uri)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Failed to get project", err.Error())
	}

	if project != nil {
		// Build HTML response
		htmlResponse := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <meta name="description" content="%s">
    <meta name="keywords" content="%s">
    <link rel="canonical" href="%s">
</head>
<body>
    <h1>%s</h1><br><br>
    <p>%s</p><br>
    <p><strong>Keywords:</strong> %s</p><br>
    <p><strong>Canonical URL:</strong> %s</p>
</body>
</html>`,
			project.MetaInfo.Title,
			project.MetaInfo.Description,
			project.MetaInfo.Keywords,
			canonicalURL,
			project.MetaInfo.Title,
			strings.ReplaceAll(project.MetaInfo.Description, "\n", "<br>"),
			project.MetaInfo.Keywords,
			canonicalURL,
		)

		return &imhttp.Response{
			Data:       htmlResponse,
			StatusCode: http.StatusOK,
		}, nil
	}

	// If project is nil, return empty project
	return &imhttp.Response{
		Data:       project,
		StatusCode: http.StatusOK,
	}, nil
}
