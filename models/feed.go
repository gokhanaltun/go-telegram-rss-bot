package models

import "gorm.io/gorm"

type Feed struct {
	gorm.Model
	Name         string
	Url          string
	LastRead     string
	Notification bool
}
