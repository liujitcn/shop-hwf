package models

import "time"

type Migration struct {
	Version   string    `gorm:"primaryKey"`
	ApplyTime time.Time `gorm:"autoCreateTime"`
}

func (Migration) TableName() string {
	return "base_migration"
}
