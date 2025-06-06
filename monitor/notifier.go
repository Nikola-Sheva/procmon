package monitor

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatID int64 = 6470119229 // ← замени на свой реальный chat_id

func InitBot(token string) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота: %v", err)
	}
}

func SendAlert(pID int32, name string, cpu float64, mem float32) {
	if bot == nil {
		return
	}
	text := fmt.Sprintf("[⚠️] PID: %d\nИмя: %s\nCPU: %.2f%%\nRAM: %.2f%%", pID, name, cpu, mem)
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке в Telegram: %v", err)
	}
}
