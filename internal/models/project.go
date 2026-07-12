package models

import "time"

type ProjectConfig struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DomainConfig struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	DomainName    string    `json:"domainName"`
	RedirectTo    string    `json:"redirectTo,omitempty"`
	SSLCertStatus string    `json:"sslCertStatus"`
	PathPrefix    string    `json:"pathPrefix"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type ServerlessFunctionCode struct {
	ID          string    `json:"id"`
	ServiceID   string    `json:"serviceId"`
	Runtime     string    `json:"runtime"`
	CodeContent string    `json:"codeContent"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EnvironmentConfig struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CanvasSummary struct {
	ID                 string             `json:"id"`
	WorkspaceID        string             `json:"workspaceId,omitempty"`
	Name               string             `json:"name"`
	Description        string             `json:"description,omitempty"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	EnvironmentsCount  int                `json:"environmentsCount"`
	AppsCount          int                `json:"appsCount"`
	DatabasesCount     int                `json:"databasesCount"`
	StorageCount       int                `json:"storageCount"`
	OnlineServices     int                `json:"onlineServices"`
	TotalServices      int                `json:"totalServices"`
	ServiceIcons       []string           `json:"serviceIcons"`
	DefaultEnvironment *EnvironmentConfig `json:"defaultEnvironment,omitempty"`
}

type EnvironmentCanvas struct {
	Environment *EnvironmentConfig `json:"environment"`
	Apps        []*AppService      `json:"apps"`
	Databases   []*Database        `json:"databases"`
	Storage     []*Storage         `json:"storage"`
}

type CreateProjectRequest struct {
	ID                 string `json:"id"`
	WorkspaceID        string `json:"workspaceId,omitempty"`
	Name               string `json:"name"`
	Description        string `json:"description,omitempty"`
	RepositoryURL      string `json:"repositoryUrl,omitempty"`
	RepositoryURLSnake string `json:"repository_url,omitempty"`
	Branch             string `json:"branch,omitempty"`
	InternalPort       int    `json:"internalPort,omitempty"`
	InternalPortSnake  int    `json:"internal_port,omitempty"`
	Domain             string `json:"domain,omitempty"`
}

type (
	SetEnvVarsRequest map[string]string
	VarsRequest       map[string]string
	Webhook           struct {
		ID                    string    `json:"id"`
		ProjectID             string    `json:"projectId"`
		URL                   string    `json:"url"`
		EventTypes            []string  `json:"eventTypes"`
		IncludePREnvironments bool      `json:"includePrEnvironments"`
		CreatedAt             time.Time `json:"createdAt"`
		UpdatedAt             time.Time `json:"updatedAt"`
	}
)

type ProjectToken struct {
	ID            string     `json:"id"`
	ProjectID     string     `json:"projectId"`
	EnvironmentID string     `json:"environmentId"`
	Name          string     `json:"name"`
	TokenPrefix   string     `json:"tokenPrefix"`
	Scopes        []string   `json:"scopes"`
	IPAllowlist   []string   `json:"ipAllowlist"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type ProjectMember struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"projectId"`
	UserID     string    `json:"userId,omitempty"`
	Email      string    `json:"email"`
	Permission string    `json:"permission"`
	Status     string    `json:"status"`
	InvitedAt  time.Time `json:"invitedAt"`
	AcceptedAt time.Time `json:"acceptedAt,omitempty"`
}

type CreateWebhookRequest struct {
	URL                   string   `json:"url"`
	EventTypes            []string `json:"eventTypes"`
	IncludePREnvironments bool     `json:"includePrEnvironments"`
}

type CreateTokenRequest struct {
	Name          string     `json:"name"`
	EnvironmentID string     `json:"environmentId"`
	Scopes        []string   `json:"scopes"`
	IPAllowlist   []string   `json:"ipAllowlist,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
}

type AddMemberRequest struct {
	Email      string `json:"email"`
	Permission string `json:"permission"`
}
