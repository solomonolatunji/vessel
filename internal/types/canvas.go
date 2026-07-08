package types

type ProjectCanvasSummary struct {
	ProjectConfig
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
	Environment *EnvironmentConfig  `json:"environment"`
	Apps        []*AppServiceConfig `json:"apps"`
	Databases   []*DatabaseConfig   `json:"databases"`
	Storage     []*StorageConfig    `json:"storage"`
}
