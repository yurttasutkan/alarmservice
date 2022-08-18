package config

// Config defines the configuration structure.
type Config struct {
	General struct {
		LogLevel                  int    `mapstructure:"log_level"`
		LogToSyslog               bool   `mapstructure:"log_to_syslog"`
		GRPCDefaultResolverScheme string `mapstructure:"grpc_default_resolver_scheme"`
	} `mapstructure:"general"`

	PostgreSQL struct {
		DSN                string `mapstructure:"dsn"`
		Automigrate        bool `mapstructure:"automigrate"`
		MaxOpenConnections int `mapstructure:"max_open_connections"`
		MaxIdleConnections int `mapstructure:"max_idle_connections"`
	} `mapstructure:"postgresql"`

	AlarmServer struct{
		API struct{
			Bind string `mapstructure:"bind"`
		} `mapstructure:"api"`
		Address string `mapstructure:"als_addr"`

	} `mapstructure:"alarm_server"`
}
	// C holds the global configuration.
	var C Config

	// Get returns the configuration.
	func Get() *Config {
		return &C
	}
