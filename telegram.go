package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
)

func InitBot(config Config) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	SetCommandList(bot)

	return bot
}

func ProcessMessages(bot *tgbotapi.BotAPI, db *gorm.DB, sensor string, config Config) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			RegisterUser(update.Message, bot, db, config)
		case "get_measures":
			if CheckUserRegistered(update.Message.Chat.ID, db) {
				GetCurrentTemp(bot, sensor)
			}
		case "set_air_border_temp":
			if CheckUserRegistered(update.Message.Chat.ID, db) {
				SetAirBorderTemp(bot, db, update.Message)
			}
		case "set_water_border_temp":
			if CheckUserRegistered(update.Message.Chat.ID, db) {
				SetWaterBorderTemp(bot, db, update.Message)
			}
		case "get_temp_borders":
			if CheckUserRegistered(update.Message.Chat.ID, db) {
				GetMeasurementSettings(bot, db)
			}
		case "set_heater_temp":
			if CheckUserRegistered(update.Message.Chat.ID, db) {
				SetHeaterTemp(bot, db, update.Message)
			}
		}
	}
}

func RegisterUser(message *tgbotapi.Message, bot *tgbotapi.BotAPI, db *gorm.DB, config Config) {
	log.Println("Creating user...")
	if message.Chat.UserName == config.Telegram.Username {
		db.FirstOrCreate(&User{ChatID: message.Chat.ID, UserName: message.Chat.UserName})
		msg := tgbotapi.NewMessage(message.Chat.ID, "Вы зарегистрированы!")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Вы не имеете доступа к этому боту!")
		bot.Send(msg)
	}
}

func CheckUserRegistered(ChatID int64, db *gorm.DB) bool {
	user := User{ChatID: ChatID}
	err := db.First(&user).Error
	if err == nil {
		return true
	}

	return false
}

func SendCriticalMessage(bot *tgbotapi.BotAPI, db *gorm.DB, air_temperature float32, water_temperature float32) {
	message := ""

	if air_temperature < float32(GetAirBorderTemperature(db)) {
		critical_air_message := fmt.Sprintf("Температура воздуха ниже критического значения! Сейчас: %.2f", air_temperature)
		message = fmt.Sprintf("%s\n%s", message, critical_air_message)
	}

	if water_temperature < float32(GetWaterBorderTemperature(db)) {
		critical_water_message := fmt.Sprintf("Температура отопления ниже критического значения! Сейчас: %.2f", water_temperature)
		message = fmt.Sprintf("%s\n%s", message, critical_water_message)
	}

	if len(message) > 0 {
		SendMessage(bot, message)
	}
}

func GetCurrentTemp(bot *tgbotapi.BotAPI, sensor string) {
	air_temp := GetTemp()
	water_temp := GetWaterTemp(sensor)
	message := fmt.Sprintf("Воздух: %.2f\nВода: %.2f", air_temp, water_temp)

	SendMessage(bot, message)
}

func SetAirBorderTemp(bot *tgbotapi.BotAPI, db *gorm.DB, message *tgbotapi.Message) {
	SettingSet(db, AirBorderTemp, message.CommandArguments())
	GetMeasurementSettings(bot, db)
}

func SetWaterBorderTemp(bot *tgbotapi.BotAPI, db *gorm.DB, message *tgbotapi.Message) {
	SettingSet(db, WaterBorderTemp, message.CommandArguments())
	GetMeasurementSettings(bot, db)
}

func GetMeasurementSettings(bot *tgbotapi.BotAPI, db *gorm.DB) {
	settings_string := fmt.Sprintf("Текущие настройки:\nВоздух - %d\nВода - %d\n",
		GetAirBorderTemperature(db),
		GetWaterBorderTemperature(db))
	SendMessage(bot, settings_string)
}

func SetHeaterTemp(bot *tgbotapi.BotAPI, db *gorm.DB, message *tgbotapi.Message) {
	heater_temp := message.CommandArguments()
	SettingSet(db, "heater_temp", heater_temp)
	SendMessage(bot, fmt.Sprintf("Температура котла установлена на %s градусов", heater_temp))
}

func SendMessage(bot *tgbotapi.BotAPI, message string) {
	msg := tgbotapi.NewMessage(253035365, message)
	bot.Send(msg)
}

func SetCommandList(bot *tgbotapi.BotAPI) {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Регистрация"},
		{Command: "get_measures", Description: "Получить текущие значения"},
		{Command: "get_temp_borders", Description: "Получить значения для оповещения"},
		{Command: "set_air_border_temp", Description: "Задать критическое значение температуры воздуха"},
		{Command: "set_water_border_temp", Description: "Задать критическое значение температуры воды"},
		{Command: "set_heater_temp", Description: "Задать температуру котла"},
	}

	config := tgbotapi.NewSetMyCommands(commands...)
	bot.Send(config)
}
