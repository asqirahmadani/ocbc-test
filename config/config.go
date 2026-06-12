package config

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name           string
	Env            string
	Port           string
	Profiler       bool
	RequestTimeout int
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
}

type DatabaseConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string
	MaxIdleConn int
}

type JWTConfig struct {
	SecretKey      string
	ExpireDuration int
}
