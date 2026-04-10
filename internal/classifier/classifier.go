package classifier

import (
	"strings"

	"github.com/lucas/netmap/internal/models"
)

// ClassifyEndpoint takes a path and returns its EndpointType based on keywords.
func ClassifyEndpoint(path string) models.EndpointType {
	lowerPath := strings.ToLower(path)

	if strings.Contains(lowerPath, "login") || strings.Contains(lowerPath, "auth") || strings.Contains(lowerPath, "signin") {
		return models.TypeAuth
	}
	if strings.Contains(lowerPath, "admin") || strings.Contains(lowerPath, "dashboard") || strings.Contains(lowerPath, "panel") {
		return models.TypeAdmin
	}
	if strings.Contains(lowerPath, "api") || strings.Contains(lowerPath, "v1") || strings.Contains(lowerPath, "v2") || strings.Contains(lowerPath, "graphql") {
		return models.TypeAPI
	}
	if path == "/" || strings.Contains(lowerPath, "home") || strings.Contains(lowerPath, "about") {
		return models.TypeGeneral
	}

	return models.TypeUnknown
}
