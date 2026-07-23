package services

import (
	"strings"

	"codedock.dev/codedock/internal/models"
	"gopkg.in/yaml.v3"
)

type ComposeParserService struct{}

func NewComposeParserService() *ComposeParserService {
	return &ComposeParserService{}
}

type ParsedComposeResult struct {
	AppServices []models.CreateAppServiceRequest `json:"appServices"`
	Databases   []models.CreateDatabaseRequest   `json:"databases"`
}

func (s *ComposeParserService) Parse(composeData []byte, projectID string) (*ParsedComposeResult, error) {
	var compose models.UserComposeFile
	if err := yaml.Unmarshal(composeData, &compose); err != nil {
		return nil, err
	}

	result := &ParsedComposeResult{
		AppServices: make([]models.CreateAppServiceRequest, 0),
		Databases:   make([]models.CreateDatabaseRequest, 0),
	}

	for name, svc := range compose.Services {
		imageLower := strings.ToLower(svc.Image)
		engine := detectDatabaseEngine(imageLower)

		if engine != "" {
			dbReq := models.CreateDatabaseRequest{
				ProjectID: projectID,
				Name:      name,
				Engine:    engine,
				Version:   extractVersion(imageLower),
				Port:      getDefaultPort(engine),
			}
			result.Databases = append(result.Databases, dbReq)
		} else {
			appReq := models.CreateAppServiceRequest{
				ProjectID:   projectID,
				Name:        name,
				ImageRef:    svc.Image,
				RuntimeMode: models.RuntimeModeWeb,
			}
			if svc.Build != nil {
				appReq.BuildEngine = string(models.BuildEngineDockerfile)
			}
			result.AppServices = append(result.AppServices, appReq)
		}
	}

	return result, nil
}

func detectDatabaseEngine(image string) models.DatabaseEngine {
	if strings.Contains(image, "postgres") {
		return models.DatabaseEnginePostgres
	}
	if strings.Contains(image, "mysql") {
		return models.DatabaseEngineMySQL
	}
	if strings.Contains(image, "redis") {
		return models.DatabaseEngineRedis
	}
	if strings.Contains(image, "mongo") {
		return models.DatabaseEngineMongoDB
	}
	if strings.Contains(image, "mariadb") {
		return models.DatabaseEngineMariaDB
	}
	if strings.Contains(image, "clickhouse") {
		return models.DatabaseEngineClickhouse
	}
	return ""
}

func extractVersion(image string) string {
	parts := strings.Split(image, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return "latest"
}

func getDefaultPort(engine models.DatabaseEngine) int {
	switch engine {
	case models.DatabaseEnginePostgres:
		return 5432
	case models.DatabaseEngineMySQL:
		return 3306
	case models.DatabaseEngineRedis:
		return 6379
	case models.DatabaseEngineMongoDB:
		return 27017
	case models.DatabaseEngineMariaDB:
		return 3306
	case models.DatabaseEngineClickhouse:
		return 9000
	default:
		return 0
	}
}
