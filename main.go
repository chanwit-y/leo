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
	factory.GenGormFile()
	// factory.TestToGrom("TRIP")

	// var tripItems schemax.TripItems
	// db.Debug().Preload("Flight").Find(&tripItems, "TAI_ID = ?", 3874)
	// db.Debug().Preload(clause.Associations).Find(&tripItems, "TAI_ID = ?", 3876)
	// var trip schemax.Trip
	// db.Debug().Preload(clause.Associations).Find(&trip, "TA_ID = ?", 1950)
	// fmt.Println(trip.TripItems[0].TaiId)

	// taiId := trip.TripItems[0].TaiId
	// var cars []schemax.Car
	// db.Debug().Preload(clause.Associations).Find(&cars, "TAI_ID = ?", taiId)

	// fmt.Println(cars[0].CarId)
}
