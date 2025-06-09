package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"procmon/monitor"
	"strconv"
	"time"
)

func main() {
	cpuThreshold := flag.Float64("cpu", 30.0, "Порог CPU (%)")
	memThreshold := flag.Float64("mem", 10.0, "Порог памяти (%)")
	interval := flag.Int("interval", 10, "Интервал мониторинга (сек)")
	flag.Parse()

	monitor.InitBot("7941396107:AAESmpNfy5YVDQ7zelunO-P9UOCBTHpmrDw")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n== Меню ==")
		fmt.Println("1) Показать авторизованных пользователей")
		fmt.Println("2) Добавить пользователя")
		fmt.Println("3) Удалить пользователя")
		fmt.Println("4) Начать мониторинг")
		fmt.Println("5) Выйти")
		fmt.Print("Выберите опцию: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			users := monitor.GetUsers()
			if len(users) == 0 {
				fmt.Println("Пользователи не найдены")
			} else {
				fmt.Println("Авторизованные пользователи:")
				for username, chatID := range users {
					fmt.Printf(" - %s (chatID: %d)\n", username, chatID)
				}
			}

		case "2":
			monitor.WaitForCodeInput()

		case "3":
			users := monitor.ListUsers()
			if len(users) == 0 {
				fmt.Println("⚠ Нет авторизованных пользователей.")
				break
			}

			// Выводим список пользователей с номерами
			fmt.Println("Список пользователей:")
			for i, user := range users {
				fmt.Printf("%d) %s\n", i+1, user)
			}

			fmt.Print("Введите username или номер для удаления: ")
			scanner.Scan()
			input := scanner.Text()

			var username string
			index, err := strconv.Atoi(input)
			if err == nil && index >= 1 && index <= len(users) {
				username = users[index-1] // преобразование номера в username
			} else {
				username = input // если это не число — используем как имя
			}

			// удаляем пользователя
			if monitor.RemoveUser(username) {
				fmt.Println("✅ Пользователь удалён:", username)
			} else {
				fmt.Println("❌ Пользователь не найден:", username)
			}

		case "4":
			users := monitor.GetUsers()
			if len(users) == 0 {
				fmt.Println("Нет авторизованных пользователей. Добавьте хотя бы одного.")
				continue
			}

			chatIDs := []int64{}
			for _, id := range users {
				chatIDs = append(chatIDs, id)
			}

			log.Println("Запуск мониторинга... Нажмите Enter для остановки.")

			stopChan := make(chan struct{}) // канал для остановки

			// Запуск мониторинга в отдельной горутине
			go func() {
				for {
					select {
					case <-stopChan:
						log.Println("Мониторинг остановлен.")
						return
					default:
						monitor.CheckProcesses(chatIDs, *cpuThreshold, *memThreshold)
						time.Sleep(time.Duration(*interval) * time.Second)
					}
				}
			}()

			// Ожидание нажатия Enter
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			// Останавливаем мониторинг
			stopChan <- struct{}{}

		case "5":
			fmt.Println("Выход...")
			return

		default:
			fmt.Println("Некорректный выбор")
		}
	}
}
