package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	DBName string `mapstructure:"POSTGRES_DB"`
	DBPort string `mapstructure:"POSTGRES_PORT"`
	DBUser string `mapstructure:"POSTGRES_USER"`
	DBHost string `mapstructure:"POSTGRES_HOST"`
	DBPass string `mapstructure:"POSTGRES_PASSWORD"`

	StorageRootUser     string `mapstructure:"MINIO_ROOT_USER"`
	StorageRootPassword string `mapstructure:"MINIO_ROOT_PASSWORD"`
	StoragePortAPI      string `mapstructure:"MINIO_PORT_API"`
	StoragePortUI       string `mapstructure:"MINIO_PORT_UI"`
	StorageBucket       string `mapstructure:"MINIO_BUCKET"`
	StorageHost         string `mapstructure:"MINIO_HOST"`
	StorageRegion       string `mapstructure:"MINIO_REGION"`

	AppURL  string `mapstructure:"APP_URL"`
	AppPort string `mapstructure:"APP_PORT"`

	Timezone  string `mapstructure:"TIMEZONE"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to load config file")
		return err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to load config file")
		return err
	}

	AppConfig = &cfg
	return nil
}
