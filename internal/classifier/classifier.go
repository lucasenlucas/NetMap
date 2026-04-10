package classifier

import (
	"strings"

	"github.com/lucas/netmap/internal/models"
)

// ClassifyEndpoint takes a path and returns its EndpointType based on keywords.
func ClassifyEndpoint(path string) models.EndpointType {
	lowerPath := strings.ToLower(path)

	// High Interest: Dev & Config (Security relevant)
	if strings.Contains(lowerPath, ".git") || strings.Contains(lowerPath, "dockerfile") || 
	   strings.Contains(lowerPath, "makefile") || strings.Contains(lowerPath, "jenkins") ||
	   strings.Contains(lowerPath, ".gitignore") || strings.Contains(lowerPath, ".gitlab-ci") {
		return models.TypeDev
	}
	if strings.Contains(lowerPath, ".env") || strings.Contains(lowerPath, "config") || 
	   strings.Contains(lowerPath, "settings") || strings.Contains(lowerPath, ".yml") || 
	   strings.Contains(lowerPath, ".json") || strings.Contains(lowerPath, "setup") ||
	   strings.Contains(lowerPath, "install") || strings.Contains(lowerPath, "backup") {
		return models.TypeConfig
	}

	// Standard Categories
	if strings.Contains(lowerPath, "login") || strings.Contains(lowerPath, "auth") || 
	   strings.Contains(lowerPath, "signin") || strings.Contains(lowerPath, "signup") ||
	   strings.Contains(lowerPath, "register") {
		return models.TypeAuth
	}
	if strings.Contains(lowerPath, "admin") || strings.Contains(lowerPath, "dashboard") || 
	   strings.Contains(lowerPath, "panel") || strings.Contains(lowerPath, "administrator") ||
	   strings.Contains(lowerPath, "manage") || strings.Contains(lowerPath, "control") {
		return models.TypeAdmin
	}
	if strings.Contains(lowerPath, "api") || strings.Contains(lowerPath, "v1") || 
	   strings.Contains(lowerPath, "v2") || strings.Contains(lowerPath, "graphql") ||
	   strings.Contains(lowerPath, "swagger") || strings.Contains(lowerPath, "rest") {
		return models.TypeAPI
	}
	if path == "/" || strings.Contains(lowerPath, "home") || strings.Contains(lowerPath, "index") || 
	   strings.Contains(lowerPath, "about") {
		return models.TypeGeneral
	}

	return models.TypeUnknown
}
