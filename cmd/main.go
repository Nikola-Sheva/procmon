package main

import (
	"flag"            // Пакет для обработки аргументов командной строки
	"log"             // Пакет для логирования
	"procmon/monitor" // Подключение собственного пакета с логикой мониторинга
	"time"            // Пакет для работы с временем
)

func main() {
	// Параметры командной строки:
	// -cpu: порог загрузки CPU, при превышении которого процесс считается подозрительным
	// -mem: порог использования оперативной памяти
	// -interval: интервал между проверками (в секундах)
	cpuThreshold := flag.Float64("cpu", 30.0, "Порог CPU (%)")
	memThreshold := flag.Float64("mem", 10.0, "Порог памяти (%)")
	interval := flag.Int("interval", 10, "Интервал мониторинга (сек)")

	// Разбор аргументов командной строки
	flag.Parse()

	// Инициализация Telegram-бота
	monitor.InitBot("7941396107:AAESmpNfy5YVDQ7zelunO-P9UOCBTHpmrDw")

	// Ожидание, пока пользователь введёт 6-значный код, полученный через Telegram
	// Возвращается имя пользователя и его chatID
	username, chatID := monitor.WaitForCodeInput()

	log.Printf("✅ Мониторинг начат для пользователя: %s\n", username)

	// Бесконечный цикл мониторинга системы
	for {
		// Проверка всех процессов, если они превышают указанные пороги — отправка уведомлений
		monitor.CheckProcesses(chatID, *cpuThreshold, *memThreshold)

		// Задержка перед следующей проверкой
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
