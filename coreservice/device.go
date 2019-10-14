package coreservice

import (
	"github.com/jinzhu/gorm"
)

type SN struct {
	gorm.Model
	Value		string			`gorm:"unique_index"`
}

type Device struct {
	gorm.Model
	SN			string
}

type DeviceService struct {
	db			*gorm.DB
}