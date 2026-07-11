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

type AuditLog struct {
	gorm.Model
	TeamID    string
	UserID    string
	Action    string
	Resource  string
	IPAddress string
	Timestamp time.Time
}
