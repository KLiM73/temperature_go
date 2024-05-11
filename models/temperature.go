package models

import "gorm.io/gorm"

type Temperature struct {
	gorm.Model
	HeaterTemp  float32
	SensorTemp  float32
	OutsideTemp float32
}
