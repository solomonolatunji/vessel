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

func (s *RailwayService) doGraphQL(ctx context.Context, token, query string, variables map[string]interface{}) ([]byte, error) {
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

type RailwayProjectNode struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RailwayProjectsResponse struct {
	Data struct {
		Projects struct {
			Edges []struct {
				Node RailwayProjectNode `json:"node"`
			} `json:"edges"`
		} `json:"projects"`
	} `json:"data"`
}

func (s *RailwayService) ListProjects(ctx context.Context, token string) ([]RailwayProjectNode, error) {
	query := `query { projects { edges { node { id name description } } } }`
	body, err := s.doGraphQL(ctx, token, query, nil)
	if err != nil {
		return nil, err
	}

	var res RailwayProjectsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	var projects []RailwayProjectNode
	for _, edge := range res.Data.Projects.Edges {
		projects = append(projects, edge.Node)
	}

	return projects, nil
}

type RailwayEnvironmentNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RailwayServiceSource struct {
	Image string `json:"image"`
	Repo  string `json:"repo"`
}

type RailwayServiceNode struct {
	ID     string               `json:"id"`
	Name   string               `json:"name"`
	Source RailwayServiceSource `json:"source"`
}

type RailwayProjectDetails struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Environments struct {
		Edges []struct {
			Node RailwayEnvironmentNode `json:"node"`
		} `json:"edges"`
	} `json:"environments"`
	Services struct {
		Edges []struct {
			Node RailwayServiceNode `json:"node"`
		} `json:"edges"`
	} `json:"services"`
}

type RailwayProjectDetailsResponse struct {
	Data struct {
		Project RailwayProjectDetails `json:"project"`
	} `json:"data"`
}

func (s *RailwayService) GetProjectDetails(ctx context.Context, token, projectID string) (*RailwayProjectDetails, error) {
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

	var res RailwayProjectDetailsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res.Data.Project, nil
}

func (s *RailwayService) ImportProject(ctx context.Context, token, projectID string, excludeRailwayVars, recreateDatabases, importData bool) error {
	details, err := s.GetProjectDetails(ctx, token, projectID)
	if err != nil {
		return err
	}

	proj, err := s.projectService.CreateProject(ctx, details.Name, details.Description)
	if err != nil {
		return err
	}

	// Create environments
	for _, envEdge := range details.Environments.Edges {
		env, err := s.projectService.CreateEnvironment(ctx, proj.ID, envEdge.Node.Name)
		if err != nil {
			return err
		}

		// Create services in each environment
		for _, svcEdge := range details.Services.Edges {
			svcNode := svcEdge.Node

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
					ProjectID:     proj.ID,
					EnvironmentID: env.ID,
					Name:          svcNode.Name,
					Engine:        engine,
					Version:       "latest",
				}
				if _, err := s.dbService.CreateDatabase(ctx, db); err != nil {
					return err
				}
			} else {
				app := &models.AppService{
					ProjectID:     proj.ID,
					EnvironmentID: env.ID,
					Name:          svcNode.Name,
					RepositoryURL: svcNode.Source.Repo,
				}
				if _, err := s.appService.CreateAppService(ctx, app); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
