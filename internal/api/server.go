package api

import (
	"net/http"
	"strings"

	"github.com/docker/docker/client"
	"github.com/solomonolatunji/vessel/internal/orchestrator"
	"github.com/solomonolatunji/vessel/internal/proxy"
	"github.com/solomonolatunji/vessel/internal/store"
)

// Server encapsulates HTTP routing, API handler dependencies, and authentication guards for the Vessel control plane.
type Server struct {
	router          *http.ServeMux
	store           *store.Store
	deployer        *orchestrator.Deployer
	proxyManager    *proxy.ProxyManager
	dockerClient    *client.Client
	tokenService    *TokenService
	dbDeployer      *orchestrator.DatabaseDeployer
	storageDeployer *orchestrator.StorageDeployer
}

// NewServer initializes a Server wired to the database store, container orchestrator, reverse proxy, and Docker client.
func NewServer(s *store.Store, deployer *orchestrator.Deployer, proxyManager *proxy.ProxyManager, dockerClient *client.Client) *Server {
	srv := &Server{
		router:          http.NewServeMux(),
		store:           s,
		deployer:        deployer,
		proxyManager:    proxyManager,
		dockerClient:    dockerClient,
		tokenService:    NewTokenService(),
		dbDeployer:      orchestrator.NewDatabaseDeployer(dockerClient, s),
		storageDeployer: orchestrator.NewStorageDeployer(dockerClient, s),
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

		if strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/api/auth/") {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "invalid authorization token format")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
