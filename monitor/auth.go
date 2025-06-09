package monitor

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Глобальные переменные
var bot *tgbotapi.BotAPI
var users = make(map[string]int64)         // username -> chatID (зарегистрированные пользователи)
var pendingCodes = make(map[string]string) // code -> username (ожидающие подтверждения)

// InitBot инициализирует Telegram-бота и запускает слушатель обновлений
func InitBot(token string) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота: %v", err)
	}

	log.Printf("✅ Бот авторизован как %s", bot.Self.UserName)

	loadUsers() // Загружаем пользователей из файла

	// Настраиваем получение обновлений
	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	// Обрабатываем обновления в отдельной горутине
	go func() {
		for update := range updates {
			handleUpdate(update)
		}
	}()
}

// handleUpdate обрабатывает входящие сообщения от Telegram
func handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return // Пропускаем не-сообщения
	}

	username := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// Создаем клавиатуру с одной кнопкой "Start"
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Start"),
		),
	)

	if text == "/start" {
		// Отправляем клавиатуру с кнопкой
		msg := tgbotapi.NewMessage(chatID, "Нажмите кнопку 'Start' чтобы получить код для привязки.")
		msg.ReplyMarkup = replyKeyboard
		bot.Send(msg)
		return
	}

	if text == "Start" {
		// Генерация кода привязки
		code := generateCode()

		// Сохраняем ожидающего пользователя
		pendingCodes[code] = username
		users[username] = chatID // Сохраняем chatID
		saveUsers()

		// Отправляем код пользователю
		msg := tgbotapi.NewMessage(chatID, "🔐 Ваш код для привязки: "+code)
		// После отправки кода можно скрыть клавиатуру
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		bot.Send(msg)
		return
	}

}

// generateCode генерирует случайный 6-значный код
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(100000 + rand.Intn(900000))
}

// GetChatIDs возвращает список всех chatID зарегистрированных пользователей
func GetChatIDs() []int64 {
	ids := []int64{}
	for _, id := range users {
		ids = append(ids, id)
	}
	return ids
}

// SendAlertTo отправляет уведомление об опасном процессе конкретному пользователю
func SendAlertTo(chatID int64, pid int32, name string, cpu float64, mem float32) {
	msg := fmt.Sprintf("[⚠] PID: %d\nИмя: %s\nCPU: %.2f%%\nRAM: %.2f%%", pid, name, cpu, mem)
	bot.Send(tgbotapi.NewMessage(chatID, msg))
}

// WaitForCodeInput() получает от пользователя код и сверяет его с
func WaitForCodeInput() (string, int64) {
	reader := bufio.NewReader(os.Stdin)

	for attempts := 1; attempts <= 3; attempts++ {
		fmt.Print("Введите код, который вы получили в Telegram: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if username, ok := pendingCodes[input]; ok {
			chatID := users[username]
			fmt.Println("✅ Привязка успешна.")
			delete(pendingCodes, input)
			saveUsers()
			return username, chatID
		} else {
			fmt.Printf("❌ Неверный код. Попытка %d из 3.\n", attempts)
		}
	}

	fmt.Println("❌ Все попытки исчерпаны. Вы можете попробовать снова.")
	return "", 0 // Рекурсивный вызов для новой серии попыток
}
