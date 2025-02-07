package structs

import (
	"errors"
)

type Bot bool
type Person bool

// Структура пользователя (для хранения данных о пользователе)
type User struct {
	id           int64
	ContextArray context // Контекст сообщений с пользователем
}

// Структура контекста, содержащая
type context struct {
	// Индекс последнего добавленного элемента в слайс, сделано для удобства
	currentSizeIndex int
	// Размер контекста, устанавливается в конструкторе newContext()
	size int
	// Сам слайс, содержащий контекст
	array []string
}

// Конструктор для User struct
func NewUser(_id int64) *User {
	return &User{
		id:           _id,
		ContextArray: newContext(),
	}
}

// Конструктор для context struct
func newContext() context {
	result := context{
		currentSizeIndex: -1,
		// Размер контекста для пользователя (в данном случа 3 сообщения от бота, 3 от пользователя)
		// При масштабировани проекта и добавлением базы данных могут возникнуть проблемы
		size: 6,
	}

	result.array = make([]string, result.size)
	return result
}

// Добавление (сохранение) контекста пользователя и чата gpt
func (c *context) AddContext(newContext string, t any) error {
	newContextArray := make([]string, c.size)

	// Проверка на переполнение контекста (можно заменить стандартными len и capacity)
	if c.currentSizeIndex+1 >= c.size {
		j := 0
		for i := c.size / 2; i < c.size; i++ {
			newContextArray[j] = c.array[i]
			j += 1
		}

		c.currentSizeIndex = c.size/2 - 1
		c.array = newContextArray
	}

	// Определение отправителя сообщения (бот или пользователь)
	switch t.(type) {
	case Bot:
		c.array[c.currentSizeIndex+1] = " CHAT GPT SMS: " + newContext
	case Person:
		c.array[c.currentSizeIndex+1] = " USER SMS: " + newContext
	default:
		return errors.New("invalid type of message sender")
	}

	// Увелечение индекса последнего добавленного элемента
	c.currentSizeIndex += 1

	return nil
}

// Получение контекста
func (c context) GetContext() string {
	// Промпт для объяснения чат gpt контекста
	var context string = "История моих прошлых сообщений в скобках, проанализируй её перед ответом на сообщение после скобок ["
	isEmptyContext := true

	// Собираем контекст и добавляем к строке контекста
	// Можно заменить стандартным strings.Join или strings.Builder
	for _, message := range c.array {
		if message != "" {
			context += message
			isEmptyContext = false
		}
	}

	// Если контекста нет, то возвращаем пустой контекст
	if isEmptyContext {
		return ""
	}

	return context + "]"
}
