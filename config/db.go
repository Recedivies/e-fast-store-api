package config

import (
	"github.com/Roixys/e-fast-store-api/exception"
	"github.com/Roixys/e-fast-store-api/model"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewPostgres(dbSource string, env string) *gorm.DB {
	var gormLogger logger.Interface = logger.Default.LogMode(logger.Info)
	if env != "development" {
		// For other environments (e.g., production), use a silent logger
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{
		Logger: gormLogger,
	})
	exception.FatalIfNeeded(err, "cannot connect to db")

	return db
}

func RunDBMigration(db *gorm.DB) {
	// db.Migrator().DropTable(&model.User{}, &model.Category{}, &model.Product{},
	// 	&model.Cart{}, &model.PaymentEvent{}, &model.PaymentEvent{},
	// )
	db.Migrator().AutoMigrate(&model.User{}, &model.Category{}, &model.Product{},
		&model.Cart{}, &model.PaymentEvent{}, &model.PaymentEvent{},
	)

	log.Info().Msg("db migrated successfully")
}
