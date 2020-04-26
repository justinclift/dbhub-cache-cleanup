package main

import "time"

// Configuration file
type TomlConfig struct {
	Admin       AdminInfo
	Auth0       Auth0Info
	DB4S        DB4SInfo
	Environment EnvInfo
	DiskCache   DiskCacheInfo
	Event       EventProcessingInfo
	Licence     LicenceInfo
	Memcache    MemcacheInfo
	Minio       MinioInfo
	Pg          PGInfo
	Sign        SigningInfo
	Web         WebInfo
}

// Config info for the admin server
type AdminInfo struct {
	Certificate    string
	CertificateKey string `toml:"certificate_key"`
	HTTPS          bool
	Server         string
}

// Auth0 connection parameters
type Auth0Info struct {
	ClientID     string
	ClientSecret string
	Domain       string
}

// Configuration info for the DB4S end point
type DB4SInfo struct {
	CAChain        string `toml:"ca_chain"`
	Certificate    string
	CertificateKey string `toml:"certificate_key"`
	Port           int
	Server         string
}

// Disk cache info
type DiskCacheInfo struct {
	Directory string
}

// Environment info
type EnvInfo struct {
	Environment string
}

// Event processing loop
type EventProcessingInfo struct {
	Delay                     time.Duration `toml:"delay"`
	EmailQueueDir             string        `toml:"email_queue_dir"`
	EmailQueueProcessingDelay time.Duration `toml:"email_queue_processing_delay"`
}

// Path to the licence files
type LicenceInfo struct {
	LicenceDir string `toml:"licence_dir"`
}

// Memcached connection parameters
type MemcacheInfo struct {
	DefaultCacheTime    int           `toml:"default_cache_time"`
	Server              string        `toml:"server"`
	ViewCountFlushDelay time.Duration `toml:"view_count_flush_delay"`
}

// Minio connection parameters
type MinioInfo struct {
	AccessKey string `toml:"access_key"`
	HTTPS     bool
	Secret    string
	Server    string
}

// PostgreSQL connection parameters
type PGInfo struct {
	Database       string
	NumConnections int `toml:"num_connections"`
	Port           int
	Password       string
	Server         string
	SSL            bool
	Username       string
}

// Used for signing DB4S client certificates
type SigningInfo struct {
	CertDaysValid    int    `toml:"cert_days_valid"`
	Enabled          bool   `toml:"enabled"`
	IntermediateCert string `toml:"intermediate_cert"`
	IntermediateKey  string `toml:"intermediate_key"`
}

type WebInfo struct {
	BaseDir              string `toml:"base_dir"`
	BindAddress          string `toml:"bind_address"`
	Certificate          string `toml:"certificate"`
	CertificateKey       string `toml:"certificate_key"`
	RequestLog           string `toml:"request_log"`
	ServerName           string `toml:"server_name"`
	SessionStorePassword string `toml:"session_store_password"`
}
