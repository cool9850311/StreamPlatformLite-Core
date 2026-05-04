package config

type Config struct {
	Server struct {
		Port         int    `mapstructure:"port"`
		Domain       string `mapstructure:"domain"`
		HTTPS        bool   `mapstructure:"https" default:"false"`
		EnableGinLog bool   `mapstructure:"enable_gin_log" default:"true"`
		LogLevel     string `mapstructure:"log_level" default:"INFO"`
	} `mapstructure:"server"`
	Frontend struct {
		Domain string `mapstructure:"domain"`
		Port   int    `mapstructure:"port"`
	} `mapstructure:"frontend"`
	PostgreSQL struct {
		DSN               string `mapstructure:"dsn"`
		AutoMigrateSchema bool
	} `mapstructure:"postgresql"`
	JWT struct {
		SecretKey string `mapstructure:"secretKey"`
	} `mapstructure:"JWT"`
	Discord struct {
		ClientID     string `mapstructure:"clientId"`
		ClientSecret string `mapstructure:"clientSecret"`
		AdminID      string `mapstructure:"adminId"`
		GuildID      string `mapstructure:"guildId"`
	} `mapstructure:"discord"`
	Redis struct {
		URI string `mapstructure:"uri"`
	} `mapstructure:"redis"`
	RateLimit struct {
		Enabled bool `json:"enabled"`

		// Login endpoints (IP dimension)
		LoginPerMinute     int64 `json:"login_per_minute"`
		OAuthInitPerMinute int64 `json:"oauth_init_per_minute"`
		LogoutPerMinute    int64 `json:"logout_per_minute"`

		// Account management
		ChangePasswordPerHour int64 `json:"change_password_per_hour"`
	}
}
