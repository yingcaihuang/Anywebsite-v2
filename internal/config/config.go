package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	ACME     ACMEConfig     `mapstructure:"acme"`
	Security SecurityConfig `mapstructure:"security"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

type ServerConfig struct {
	Port   string `mapstructure:"port"`
	Mode   string `mapstructure:"mode"`
	Domain string `mapstructure:"domain"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

type ACMEConfig struct {
	Email   string `mapstructure:"email"`
	Staging bool   `mapstructure:"staging"`
}

type SecurityConfig struct {
	JWTSecret string   `mapstructure:"jwt_secret"`
	APIKeys   []string `mapstructure:"api_keys"`
}

type StorageConfig struct {
	StaticPath  string `mapstructure:"static_path"`
	UploadsPath string `mapstructure:"uploads_path"`
	CertsPath   string `mapstructure:"certs_path"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置环境变量前缀
	viper.SetEnvPrefix("SHS")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
