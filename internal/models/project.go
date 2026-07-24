package models

import "time"

type ProjectConfig struct {
	ID          string    `json:"id" db:"id"`
	ServerID    string    `json:"serverId,omitempty" db:"server_id"` // Node where this project runs
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type SSLCertStatus string

const (
	SSLCertStatusPending SSLCertStatus = "pending"
	SSLCertStatusIssued  SSLCertStatus = "issued"
	SSLCertStatusFailed  SSLCertStatus = "failed"
)

type MemberPermission string

const (
	MemberPermissionOwner  MemberPermission = "owner"
	MemberPermissionAdmin  MemberPermission = "admin"
	MemberPermissionMember MemberPermission = "member"
)

type MemberStatus string

const (
	MemberStatusPending  MemberStatus = "pending"
	MemberStatusAccepted MemberStatus = "accepted"
)

type DomainConfig struct {
	ID            string        `json:"id" db:"id"`
	ServiceID     string        `json:"serviceId" db:"service_id"`
	DomainName    string        `json:"domainName" db:"domain_name"`
	RedirectTo    string        `json:"redirectTo,omitempty" db:"redirect_to"`
	SSLCertStatus SSLCertStatus `json:"sslCertStatus" db:"ssl_cert_status"`
	PathPrefix    string        `json:"pathPrefix" db:"path_prefix"`
	CreatedAt     time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time     `json:"updatedAt" db:"updated_at"`
}

type ServerlessFunctionCode struct {
	ID          string    `json:"id" db:"id"`
	ServiceID   string    `json:"serviceId" db:"service_id"`
	Runtime     string    `json:"runtime" db:"runtime"`
	CodeContent string    `json:"codeContent" db:"code_content"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type EnvironmentConfig struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"projectId" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	IsDefault bool      `json:"isDefault" db:"is_default"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CanvasSummary struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Description        string             `json:"description,omitempty"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	EnvironmentsCount  int                `json:"environmentsCount"`
	AppsCount          int                `json:"appsCount"`
	DatabasesCount     int                `json:"databasesCount"`
	OnlineServices     int                `json:"onlineServices"`
	TotalServices      int                `json:"totalServices"`
	ServiceIcons       []string           `json:"serviceIcons"`
	DefaultEnvironment *EnvironmentConfig `json:"defaultEnvironment,omitempty"`
}

type EnvironmentCanvas struct {
	Environment *EnvironmentConfig `json:"environment"`
	Apps        []*AppService      `json:"apps"`
	Databases   []*Database        `json:"databases"`
}

type CreateProjectRequest struct {
	ID                 string `json:"id"`
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
)
type ProjectToken struct {
	ID            string     `json:"id" db:"id"`
	ProjectID     string     `json:"projectId" db:"project_id"`
	EnvironmentID string     `json:"environmentId" db:"environment_id"`
	Name          string     `json:"name" db:"name"`
	TokenPrefix   string     `json:"tokenPrefix" db:"token_prefix"`
	Scopes        []string   `json:"scopes" db:"-"`
	IPAllowlist   []string   `json:"ipAllowlist" db:"-"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
}

type ProjectMember struct {
	ID         string           `json:"id" db:"id"`
	ProjectID  string           `json:"projectId" db:"project_id"`
	UserID     string           `json:"userId,omitempty" db:"user_id"`
	Email      string           `json:"email" db:"email"`
	Permission MemberPermission `json:"permission" db:"permission"`
	Status     MemberStatus     `json:"status" db:"status"`
	InvitedAt  time.Time        `json:"invitedAt" db:"invited_at"`
	AcceptedAt time.Time        `json:"acceptedAt,omitempty" db:"accepted_at"`
}

type CreateTokenRequest struct {
	Name          string     `json:"name"`
	EnvironmentID string     `json:"environmentId"`
	Scopes        []string   `json:"scopes"`
	IPAllowlist   []string   `json:"ipAllowlist,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
}

type AddMemberRequest struct {
	Email      string           `json:"email"`
	Permission MemberPermission `json:"permission"`
}
