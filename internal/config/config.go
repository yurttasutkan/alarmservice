package config

// Config defines the configuration structure.
type Config struct {
	General struct {
		LogLevel                  int    `mapstructure:"log_level"`
		LogToSyslog               bool   `mapstructure:"log_to_syslog"`
		PasswordHashIterations    int    `mapstructure:"password_hash_iterations"`
		GRPCDefaultResolverScheme string `mapstructure:"grpc_default_resolver_scheme"`
	} `mapstructure:"general"`

	PostgreSQL struct {
		DSN                string `mapstructure:"dsn"`
		Automigrate        bool
		MaxOpenConnections int `mapstructure:"max_open_connections"`
		MaxIdleConnections int `mapstructure:"max_idle_connections"`
	} `mapstructure:"postgresql"`

	AlarmService struct{
		ID string `mapstructure:"id"`
		Address string `mapstructure:"als_addr"`

	} `mapstructure:"alarm_service"`
}
	// C holds the global configuration.
	var C Config

	// Get returns the configuration.
	func Get() *Config {
		return &C
	}
