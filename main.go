package main

import (
	"leo/pkg/env"
	"leo/pkg/utils"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	dsn := env.Env().CONNECTION_STRING
	db, _ := gorm.Open(sqlserver.Open(dsn))
	factory := utils.NewFactory(utils.NewMsSql(db))
	factory.CreateGorm("TRIP")
}
