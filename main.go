package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var config = InitConfig()
var ds18b20 = InitDs18b20()
var bot = InitBot(config)
var db = InitDB(config)

func main() {
	InitSettings(db)
	go ProcessMessages(bot, db, ds18b20, config)

	GreetingMessage := fmt.Sprintf("Устройство запущено! Текущие показания:\nВоздух - %.2f\n"+
		"Вода - %.2f\n\nПороговые значения:\nВоздух - %s\nВода - %s\n\nТемпература котла: %s",
		GetTemp(),
		GetWaterTemp(ds18b20),
		SettingGet(db, AirBorderTemp),
		SettingGet(db, WaterBorderTemp),
		SettingGet(db, HeaterTemp),
	)
	log.Println("Taking measures...")
	SendMessage(bot, GreetingMessage)

	for true {
		log.Println("Taking measures...")
		air_temp := GetTemp()
		water_temp := GetWaterTemp(ds18b20)
		SendCriticalMessage(bot, db, air_temp, water_temp)
		SendData(air_temp, water_temp, config)
		time.Sleep(config.Backend.SendTimeout * time.Second)
	}
}

func InitDs18b20() string {
	sensors, err := Sensors()
	if err != nil {
		SendMessage(bot, "DS18B20: Ошибка инициализации: "+err.Error())
		log.Fatal(err)
	}

	return sensors[0]
}

func GetWaterTemp(sensor string) float32 {
	water_temp, err := Temperature(sensor)
	if err != nil {
		SendMessage(bot, "DS18B20: Ошибка чтения")
		log.Fatal(err)
	}

	return float32(water_temp)
}

func SendData(temp float32, water_temp float32, config Config) {
	heater_temp, _ := strconv.ParseFloat(SettingGet(db, HeaterTemp), 32)
	reqMap := map[string]float32{
		"temperature":        temp,
		"water_temperature":  water_temp,
		"heater_temperature": float32(heater_temp),
	}

	jsonBody, _ := json.Marshal(reqMap)
	bodyReader := bytes.NewBuffer(jsonBody)
	log.Println(bodyReader)
	log.Println("sending request")
	_, err := http.Post(config.Backend.Url, "application/json", bodyReader)
	if err != nil {
		log.Fatal(err)
	}
}
