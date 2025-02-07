package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/structs"
	"net/url"

	"github.com/glossd/fetch" // для более простых http запросов без использования net/url

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // telegram bot
)

var object1 structs.Person
var object2 structs.Bot

func main() {

	defer func() {
		log.Print("[ACTION]: Bot disabled.")
	}()

	// Инициализация бота с telegram api
	bot, err := tgbotapi.NewBotAPI("YOUR_API")

	// Обработка ошибок
	if err != nil {
		log.Fatal(err)
	}

	log.Print("[ACTION]: Bot started.")

	bot.Debug = false

	// Настройка конфига
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 30

	// Канал для получения обновлений бота
	updatesChannel := bot.GetUpdatesChan(config)

	// Канал для сообщений, которые будут попадать в него из канала обновлений
	messagesChannel := make(chan tgbotapi.Update)

	// Вызов анонимной горутины для получения обновлений бота
	go func() {
		for update := range updatesChannel {
			if update.Message != nil {
				messagesChannel <- update
			}
		}

		close(messagesChannel)
	}()

	// Вызов горутины с обработкой сообщений
	go messageHandler(bot, messagesChannel)

	// Ожидание завершения
	select {}
}

// Обработка сообщений
func messageHandler(bot *tgbotapi.BotAPI, messagesChannel chan tgbotapi.Update) {
	// Временное хранилище пользователей вместо использования базы данных
	users := make(map[int64]*structs.User)

	// Обработка входящих сообщений из канала сообщений
	for update := range messagesChannel {
		// Проверяем существование пользователя в нашей мапе пользователей
		_, check := users[update.Message.Chat.ID]
		if !check {
			users[update.Message.Chat.ID] = structs.NewUser(update.Message.Chat.ID)
		}

		// Получение пользователя по указателю для внесения изменения для пользователя в мапе
		user := users[update.Message.Chat.ID]

		// Запрос к chat gpt и сохранение контекста
		result := aiRequest(user.ContextArray.GetContext() + " | Новое сообщение: [" + update.Message.Text + "]")
		// Формирование нового сообщения
		message := tgbotapi.NewMessage(update.Message.Chat.ID, result)

		// Добавление (сохранение) контекста от пользователя к нашему слайсу
		err := user.ContextArray.AddContext(update.Message.Text, object1)
		if err != nil {
			panic(err)
		}

		// Добавление (сохранение) контекста от бота к нашему слайсу
		err = user.ContextArray.AddContext(result, object2)
		if err != nil {
			panic(err)
		}

		// Отправка сообщения
		bot.Send(message)
	}
}

// http запрос к free gpt ai
func aiRequest(message string) string {

	// Формирование headers для запроса
	headers := map[string]string{"Accept": "application/json"}
	// encode для нашего текста в url формат (требуется для формирования запроса к chat gpt, подробнее на github free gpt ai)
	encodedMessage := url.QueryEscape(message)

	// Формирование итогового запроса
	url := fmt.Sprintf("https://free-unoficial-gpt4o-mini-api-g70n.onrender.com/chat/?query=%s", encodedMessage)

	// http запрос с headers и нажим сообщением
	result, err := fetch.Get[string](url, fetch.Config{Headers: headers})

	// Обработка ошибок http запроса
	if err != nil {
		log.Print(err)
		panic(err)
	}

	// Создание мапы для преобразования полученного ответа из json в string
	var responceData map[string]interface{}
	// Преобразование json -> string
	err = json.Unmarshal([]byte(result), &responceData)

	// Обработка ошибок
	if err != nil {
		panic(err)
	}

	// Возврат результата запроса, преобразованного к типу string (преобразование интерфейсов)
	return responceData["results"].(string)
}
