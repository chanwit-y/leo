package utils

import "gorm.io/gorm"

type Database struct {
	db *gorm.DB
}

func New(db *gorm.DB) Database {
	return Database{db}
}

func (db *Database) Query(sql string) any {
	return ""
}
