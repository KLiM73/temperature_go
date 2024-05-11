package main

import (
	"gorm.io/gorm"
	"strconv"
)

const HeaterTemp = "heater_temp"
const AirBorderTemp = "air_border_temp"
const WaterBorderTemp = "water_border_temp"

type Setting struct {
	gorm.Model
	Key   string
	Value string
}

func InitSettings(db *gorm.DB) {
	InitialSettings := map[string]string{
		HeaterTemp:      "60",
		AirBorderTemp:   "10",
		WaterBorderTemp: "50",
	}

	for key, value := range InitialSettings {
		setted_value := SettingGet(db, key)

		if len(setted_value) == 0 {
			SettingSet(db, key, value)
		}
	}
}

func SettingSet(db *gorm.DB, key string, value string) {
	var setting Setting
	db.FirstOrCreate(&setting, "key = ?", key)

	setting.Key = key
	setting.Value = value
	db.Save(&setting)
}

func SettingGet(db *gorm.DB, key string) string {
	var setting Setting
	db.First(&setting, "key = ?", key)

	return setting.Value
}

func GetAirBorderTemperature(db *gorm.DB) int16 {
	temp := SettingGet(db, AirBorderTemp)
	temp_int, _ := strconv.ParseInt(temp, 10, 16)
	temp_int16 := int16(temp_int)
	return temp_int16
}

func GetWaterBorderTemperature(db *gorm.DB) int16 {
	temp := SettingGet(db, WaterBorderTemp)
	temp_int, _ := strconv.ParseInt(temp, 10, 16)
	temp_int16 := int16(temp_int)
	return temp_int16
}
