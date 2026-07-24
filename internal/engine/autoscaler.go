package engine

import (
	"context"
	"fmt"
	"log"
	"time"

	"codedock.run/codedock/internal/models"
)

type AppRepository interface {
	ListAll(ctx context.Context) ([]*models.AppService, error)
	Update(ctx context.Context, app *models.AppService) error
}

type DeploymentCreator interface {
	CreateDeployment(ctx context.Context, d *models.Deployment) (*models.Deployment, error)
	ExecuteDeploymentAsync(d *models.Deployment)
}

type AutoscalerWorker struct {
	appRepo           AppRepository
	statsMonitor      *StatsMonitor
	deploymentService DeploymentCreator
	ticker            *time.Ticker
	quit              chan struct{}
}

func NewAutoscalerWorker(appRepo AppRepository, statsMonitor *StatsMonitor, deploymentService DeploymentCreator) *AutoscalerWorker {
	return &AutoscalerWorker{
		appRepo:           appRepo,
		statsMonitor:      statsMonitor,
		deploymentService: deploymentService,
		quit:              make(chan struct{}),
	}
}

func (a *AutoscalerWorker) Start() {
	a.ticker = time.NewTicker(2 * time.Minute)
	go func() {
		for {
			select {
			case <-a.ticker.C:
				a.checkAndScale(context.Background())
			case <-a.quit:
				a.ticker.Stop()
				return
			}
		}
	}()
}

func (a *AutoscalerWorker) Stop() {
	if a.quit != nil {
		close(a.quit)
	}
}

func (a *AutoscalerWorker) checkAndScale(ctx context.Context) {
	apps, err := a.appRepo.ListAll(ctx)
	if err != nil {
		log.Printf("[Autoscaler] Failed to list apps: %v\n", err)
		return
	}

	for _, app := range apps {
		if app.Status != models.AppServiceStatusRunning {
			continue
		}

		health, err := a.statsMonitor.GetHealth(ctx, app.ContainerID)
		if err != nil {
			continue
		}

		cpu := health.CPUUsagePercentage
		currentReplicas := app.Replicas
		if currentReplicas <= 0 {
			currentReplicas = 1
		}

		targetReplicas := currentReplicas
		if cpu > 80.0 && currentReplicas < 5 {
			targetReplicas++
		} else if cpu < 20.0 && currentReplicas > 1 {
			targetReplicas--
		}

		if targetReplicas != currentReplicas {
			log.Printf("[Autoscaler] Scaling %s from %d to %d replicas (CPU: %.2f%%)\n", app.Name, currentReplicas, targetReplicas, cpu)
			app.Replicas = targetReplicas
			if err := a.appRepo.Update(ctx, app); err != nil {
				log.Printf("[Autoscaler] Failed to update replicas: %v\n", err)
				continue
			}

			newDep := &models.Deployment{
				ServiceID:     app.ID,
				EnvironmentID: app.EnvironmentID,
				ProjectID:     app.ProjectID,
				Status:        "BUILDING",
				CommitMessage: fmt.Sprintf("Auto-Scaling to %d replicas", targetReplicas),
				Branch:        app.Branch,
				Trigger:       "Auto-Scaler",
			}
			created, err := a.deploymentService.CreateDeployment(ctx, newDep)
			if err == nil {
				a.deploymentService.ExecuteDeploymentAsync(created)
			}
		}
	}
}
