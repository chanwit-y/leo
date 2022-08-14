package utils

import "gorm.io/gorm"

type MsSql struct {
	db *gorm.DB
}

func NewMsSql(db *gorm.DB) MsSql {
	return MsSql{db}
}

func (ms *MsSql) Query(sql string) *gorm.DB {
	return ms.db.Raw(sql)
}
