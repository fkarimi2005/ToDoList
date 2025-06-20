package model

type Configs struct {
	AuthParams AuthParams `json:"auth_params"`
	LogParams  LogParams  `json:"log_params"`
	AppParams  AppParams  `json:"app_params"`
	Postgres   Postgres   `json:"postgres"`
}
type AuthParams struct {
	JwtSecretKey  string `json:"jwt_secret_key"`
	JwtTtlMinutes int    `json:"jwt_ttl_minutes"`
}
type LogParams struct {
	LogDirectory     string `json:"log_directory"`
	LogInfo          string `json:"log_info"`
	LogError         string `json:"log_error"`
	LogWarn          string `json:"log_warn"`
	LogDebug         string `json:"log_debug"`
	MaxSizeMegabytes int    `json:"max_size_megabytes"`
	MaxBackups       int    `json:"max_backups"`
	MaxAgeDays       int    `json:"max_age_days"`
	Compress         bool   `json:"compress"`
	LocalTime        bool   `json:"local_time"`
}
type AppParams struct {
	ServerURL  string `json:"server_url"`
	ServerName string `json:"server_name"`
	AppVersion string `json:"app_version"`
	PortRun    string `json:"port_run"`
	GinMode    string `json:"gin_mode"`
}
type Postgres struct {
	User     string `json:"username"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}
