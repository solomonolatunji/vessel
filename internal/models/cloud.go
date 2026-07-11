package models

import (
	"time"

	"gorm.io/gorm"
)

type CloudTeam struct {
	gorm.Model
	Name             string `gorm:"uniqueIndex"`
	Plan             string `gorm:"default:'hobby'"`
	StripeCustomerID string
	PaddleCustomerID string `gorm:"index"`
	CustomDomain     string
	LogoURL          string
	PrimaryColor     string
}

type CloudServer struct {
	gorm.Model
	TeamID    uint
	ServerID  string `gorm:"uniqueIndex"`
	Name      string
	IPAddress string
	IsActive  bool
	LastPing  time.Time
}

type CloudUsageLog struct {
	gorm.Model
	TeamID         uint
	Deployments    int
	ContainerHours int
	BandwidthGB    int
	ReportedAt     time.Time
}

type CloudTelemetryLog struct {
	gorm.Model
	InstanceID    string `gorm:"index"`
	Version       string
	OS            string
	Arch          string
	ActiveServers int
	ActiveApps    int
	ReportedAt    time.Time
}

// CloudUser represents a user of the Vessel Cloud platform.
// It is backed by the cloud_users PostgreSQL table (managed via database/sql, not gorm).
type CloudUser struct {
	ID                   string
	Email                string
	FullName             string
	PasswordHash         string
	Role                 string // 'user' | 'admin' | 'staff'
	EmailVerified        bool
	VerifiedAt           *time.Time
	VerifyToken          string
	VerifyTokenExpiresAt *time.Time
	OTPCode              string
	OTPExpiresAt         *time.Time
	CreatedAt            time.Time
}
