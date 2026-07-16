package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"vessl.dev/vessl/internal/models"
)

type RailwayService struct {
	client         *http.Client
	projectService *ProjectService
	envService     *EnvironmentService
	appService     *AppService
	dbService      *DatabaseService
}

func NewRailwayService(ps *ProjectService, es *EnvironmentService, as *AppService, ds *DatabaseService) *RailwayService {
	return &RailwayService{
		client:         &http.Client{},
		projectService: ps,
		envService:     es,
		appService:     as,
		dbService:      ds,
	}
}

func (s *RailwayService) doGraphQL(ctx context.Context, token, query string, variables map[string]any) ([]byte, error) {
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://backboard.railway.com/graphql/v2", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("railway api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

func (s *RailwayService) ListProjects(ctx context.Context, token string) ([]models.RailwayProject, error) {
	query := `query { projects { edges { node { id name description } } } }`
	body, err := s.doGraphQL(ctx, token, query, nil)
	if err != nil {
		return nil, err
	}

	var res models.RailwayProjectsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	var projects []models.RailwayProject
	for _, edge := range res.Data.Projects.Edges {
		projects = append(projects, edge.Node)
	}

	return projects, nil
}

func (s *RailwayService) GetProjectDetails(ctx context.Context, token, projectID string) (*models.RailwayProject, error) {
	query := `query($id: String!) {
		project(id: $id) {
			id
			name
			description
			environments {
				edges {
					node {
						id
						name
					}
				}
			}
			services {
				edges {
					node {
						id
						name
						source {
							image
							repo
						}
					}
				}
			}
		}
	}`
	variables := map[string]interface{}{"id": projectID}
	body, err := s.doGraphQL(ctx, token, query, variables)
	if err != nil {
		return nil, err
	}

	var res models.RailwayProjectDetailsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res.Data.Project, nil
}

type RailwayImportOptions struct {
	Token              string
	ProjectID          string
	ExcludeRailwayVars bool
	RecreateDatabases  bool
	ImportData         bool
}

func (s *RailwayService) ImportProject(ctx context.Context, opts RailwayImportOptions) error {
	details, err := s.GetProjectDetails(ctx, opts.Token, opts.ProjectID)
	if err != nil {
		return err
	}

	proj, err := s.projectService.CreateProject(ctx, details.Name, details.Description)
	if err != nil {
		return err
	}

	for _, envEdge := range details.Environments.Edges {
		if err := s.importEnvironment(ctx, proj.ID, envEdge.Node.Name, details.Services.Edges); err != nil {
			return err
		}
	}

	return nil
}

func (s *RailwayService) importEnvironment(ctx context.Context, projectID, envName string, services []models.RailwayServiceEdge) error {
	env, err := s.projectService.CreateEnvironment(ctx, projectID, envName)
	if err != nil {
		return err
	}

	for _, svcEdge := range services {
		if err := s.importService(ctx, projectID, env.ID, svcEdge.Node); err != nil {
			return err
		}
	}
	return nil
}

func (s *RailwayService) importService(ctx context.Context, projectID, envID string, svcNode models.RailwayService) error {
	isDB := false
	engine := ""
	image := strings.ToLower(svcNode.Source.Image)

	if strings.Contains(image, "postgres") {
		isDB = true
		engine = "postgres"
	} else if strings.Contains(image, "redis") {
		isDB = true
		engine = "redis"
	} else if strings.Contains(image, "mysql") {
		isDB = true
		engine = "mysql"
	}

	if isDB {
		db := &models.Database{
			ProjectID:     projectID,
			EnvironmentID: envID,
			Name:          svcNode.Name,
			Engine:        models.DatabaseEngine(engine),
			Version:       "latest",
		}
		_, err := s.dbService.CreateDatabase(ctx, db)
		return err
	}

	app := &models.AppService{
		ProjectID:     projectID,
		EnvironmentID: envID,
		Name:          svcNode.Name,
		RepositoryURL: svcNode.Source.Repo,
	}
	_, err := s.appService.CreateAppService(ctx, app)
	return err
}
