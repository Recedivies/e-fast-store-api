package main

import (
	"github.com/Roixys/e-fast-store-api/api"
	"github.com/Roixys/e-fast-store-api/config"
)

func main() {
	configuration := config.LoadConfig(".")

	db := config.NewPostgres(configuration.DBSource, configuration.Environment)

	config.RunDBMigration(db)

	api.RunGinServer(configuration, db)
}
