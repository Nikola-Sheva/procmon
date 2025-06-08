package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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

	if text == "/start" {
		// Генерация кода привязки
		code := generateCode()

		// Сохраняем ожидающего пользователя
		pendingCodes[code] = username
		users[username] = chatID // Сохраняем chatID
		saveUsers()

		// Отправляем код пользователю
		msg := tgbotapi.NewMessage(chatID, "🔐 Ваш код для привязки: "+code)
		bot.Send(msg)
	}

}

// generateCode генерирует случайный 6-значный код
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(100000 + rand.Intn(900000))
}

// loadUsers загружает пользователей из файла users.json
func loadUsers() {
	data, err := os.ReadFile("users.json")
	if err == nil {
		json.Unmarshal(data, &users)
	}
}

// saveUsers сохраняет текущих пользователей в файл users.json
func saveUsers() {
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile("users.json", data, 0644)
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
	var input string
	fmt.Print("Введите код, который вы получили в Telegram: ")
	fmt.Scanln(&input)

	for username, chatID := range users {
		if code, ok := pendingCodes[input]; ok && code == username {
			fmt.Println("Привязка успешна.")
			delete(pendingCodes, input)
			return username, chatID
		}
	}

	fmt.Println("❌ Неверный код.")
	os.Exit(1)
	return "", 0
}
