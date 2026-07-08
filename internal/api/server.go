package api

import (
	"net/http"
	"strings"

	"github.com/docker/docker/client"
	"github.com/solomonolatunji/vessel/internal/orchestrator"
	"github.com/solomonolatunji/vessel/internal/proxy"
	"github.com/solomonolatunji/vessel/internal/services"
	"github.com/solomonolatunji/vessel/internal/store"
)

// Server encapsulates HTTP routing, API handler dependencies, and authentication guards for the Vessel control plane.
type Server struct {
	router          *http.ServeMux
	store           *store.Store
	deployer        *orchestrator.Deployer
	proxyManager    *proxy.ProxyManager
	dockerClient    *client.Client
	tokenService    *services.TokenService
	dbDeployer      *orchestrator.DatabaseDeployer
	storageDeployer *orchestrator.StorageDeployer
	cronManager     *orchestrator.CronManager
	cronService     *services.CronService
	serviceLinker   *services.ServiceLinker
	gitService      *services.GitService
}

// NewServer initializes a Server wired to the database store, container orchestrator, reverse proxy, and Docker client.
func NewServer(s *store.Store, deployer *orchestrator.Deployer, proxyManager *proxy.ProxyManager, dockerClient *client.Client) *Server {
	cronMgr := orchestrator.NewCronManager(dockerClient, s)
	_ = cronMgr.Start()

	srv := &Server{
		router:          http.NewServeMux(),
		store:           s,
		deployer:        deployer,
		proxyManager:    proxyManager,
		dockerClient:    dockerClient,
		tokenService:    services.NewTokenService(),
		dbDeployer:      orchestrator.NewDatabaseDeployer(dockerClient, s),
		storageDeployer: orchestrator.NewStorageDeployer(dockerClient, s),
		cronManager:     cronMgr,
		cronService:     services.NewCronService(s, cronMgr),
		serviceLinker:   services.NewServiceLinker(s),
		gitService:      services.NewGitService(s),
	}
	if srv.deployer != nil {
		srv.deployer.EnvProvider = srv.serviceLinker.GetLinkedEnvironmentVariables
	}
	srv.registerRoutes()
	return srv
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("POST /api/auth/register", s.handleRegister)
	s.router.HandleFunc("POST /api/auth/login", s.handleLogin)
	s.router.HandleFunc("GET /api/auth/me", s.RequireAuth(s.handleGetCurrentUser))
	s.router.HandleFunc("POST /api/auth/logout", s.handleLogout)

	s.router.HandleFunc("GET /api/projects", s.handleListProjects)
	s.router.HandleFunc("POST /api/projects", s.handleCreateProject)
	s.router.HandleFunc("GET /api/projects/{id}", s.handleGetProject)
	s.router.HandleFunc("DELETE /api/projects/{id}", s.handleDeleteProject)
	s.router.HandleFunc("POST /api/projects/{id}/deploy", s.handleDeployProject)

	s.router.HandleFunc("GET /api/projects/{id}/domains", s.handleListDomains)
	s.router.HandleFunc("POST /api/projects/{id}/domains", s.handleAddDomain)
	s.router.HandleFunc("DELETE /api/domains/{id}", s.handleDeleteDomain)

	s.router.HandleFunc("GET /api/projects/{id}/env", s.handleGetEnvVars)
	s.router.HandleFunc("PUT /api/projects/{id}/env", s.handleSetEnvVars)

	s.router.HandleFunc("GET /api/databases", s.RequireAuth(s.handleListDatabases))
	s.router.HandleFunc("POST /api/databases", s.RequireAuth(s.handleCreateDatabase))
	s.router.HandleFunc("GET /api/databases/{id}", s.RequireAuth(s.handleGetDatabase))
	s.router.HandleFunc("DELETE /api/databases/{id}", s.RequireAuth(s.handleDeleteDatabase))
	s.router.HandleFunc("POST /api/databases/{id}/start", s.RequireAuth(s.handleStartDatabase))
	s.router.HandleFunc("POST /api/databases/{id}/stop", s.RequireAuth(s.handleStopDatabase))

	s.router.HandleFunc("GET /api/storage", s.RequireAuth(s.handleListStorage))
	s.router.HandleFunc("POST /api/storage", s.RequireAuth(s.handleCreateStorage))
	s.router.HandleFunc("GET /api/storage/{id}", s.RequireAuth(s.handleGetStorage))
	s.router.HandleFunc("DELETE /api/storage/{id}", s.RequireAuth(s.handleDeleteStorage))
	s.router.HandleFunc("POST /api/storage/{id}/start", s.RequireAuth(s.handleStartStorage))
	s.router.HandleFunc("POST /api/storage/{id}/stop", s.RequireAuth(s.handleStopStorage))

	s.router.HandleFunc("GET /api/jobs", s.RequireAuth(s.handleJobs))
	s.router.HandleFunc("POST /api/jobs", s.RequireAuth(s.handleJobs))
	s.router.HandleFunc("GET /api/jobs/{id}", s.RequireAuth(s.handleJobDetail))
	s.router.HandleFunc("DELETE /api/jobs/{id}", s.RequireAuth(s.handleJobDetail))
	s.router.HandleFunc("POST /api/jobs/{id}/trigger", s.RequireAuth(s.handleJobDetail))

	s.router.HandleFunc("POST /api/git/connect", s.RequireAuth(s.handleConnectGitProvider))
	s.router.HandleFunc("GET /api/git/status", s.RequireAuth(s.handleGetGitProvidersStatus))
	s.router.HandleFunc("DELETE /api/git/connect/{provider}", s.RequireAuth(s.handleDisconnectGitProvider))
	s.router.HandleFunc("GET /api/git/repos", s.RequireAuth(s.handleListGitRepositories))
	s.router.HandleFunc("POST /api/webhooks/git/{projectId}", s.handleGitWebhook)
	s.router.HandleFunc("GET /api/canvas/projects", s.RequireAuth(s.ListProjectCanvasSummaries))
	s.router.HandleFunc("GET /api/projects/{id}/summary", s.RequireAuth(s.GetProjectCanvasSummary))
	s.router.HandleFunc("GET /api/environments/{id}/canvas", s.RequireAuth(s.GetEnvironmentCanvas))

	s.router.HandleFunc("POST /api/projects/{id}/environments", s.RequireAuth(s.CreateEnvironment))
	s.router.HandleFunc("GET /api/projects/{id}/environments", s.RequireAuth(s.ListEnvironments))
	s.router.HandleFunc("DELETE /api/environments/{id}", s.RequireAuth(s.DeleteEnvironment))

	s.router.HandleFunc("POST /api/environments/{id}/apps", s.RequireAuth(s.CreateAppService))
	s.router.HandleFunc("GET /api/environments/{id}/apps", s.RequireAuth(s.ListAppServicesByEnvironment))
	s.router.HandleFunc("GET /api/apps/{id}", s.RequireAuth(s.GetAppService))
	s.router.HandleFunc("DELETE /api/apps/{id}", s.RequireAuth(s.DeleteAppService))

	s.router.HandleFunc("GET /ws/terminal/{id}", s.handleTerminalWebSocket)
}

// Handler returns the root HTTP handler wrapped with global CORS and authentication middleware.
func (s *Server) Handler() http.Handler {
	return s.corsMiddleware(s.router)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/api/auth/") && !strings.HasPrefix(r.URL.Path, "/api/webhooks/") {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "invalid authorization token format")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
