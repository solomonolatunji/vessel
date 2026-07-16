package models

type AuditLog struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Details   string `json:"details"`
	IPAddress string `json:"ipAddress"`
	CreatedAt string `json:"createdAt"`
}
