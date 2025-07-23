package database

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"starter/config"
	"starter/internal/core/role"
	"starter/internal/core/user"
)

func NewPostgresConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPass,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
		config.AppConfig.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Panic().
			Err(err).
			Str("dsn", dsn).
			Msg("unable to connect to database")

		return nil, err
	}

	err = db.AutoMigrate(&user.User{}, &role.Role{})
	if err != nil {
		log.Panic().
			Err(err).
			Msg("unable to migrate the database")
		return nil, err
	}

	return db, nil
}
