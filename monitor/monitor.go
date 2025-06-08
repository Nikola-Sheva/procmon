package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// CheckProcesses проверяет активные процессы и отправляет уведомление,
// если превышены пороги по CPU или памяти.
func CheckProcesses(chatID int64, cpuThreshold, memThreshold float64) {
	// Получаем список всех процессов в системе
	procs, err := process.Processes()
	if err != nil {
		log.Fatalf("❌ Ошибка при получении процессов: %v", err)
	}

	for _, p := range procs {
		// Получаем имя процесса
		name, _ := p.Name()

		// Получаем текущую загрузку CPU и использование памяти процессом
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()

		// Список системных процессов, которые мы игнорируем
		var whitelist = map[string]bool{
			"System":      true,
			"Idle":        true,
			"svchost.exe": true,
		}

		// Если процесс в белом списке — пропускаем
		if whitelist[name] {
			continue
		}

		// Если процесс превышает указанные пороги — формируем сообщение
		if cpu > cpuThreshold || float64(mem) > memThreshold {
			msg := fmt.Sprintf("[⚠] PID: %d | Имя: %s | CPU: %.2f%% | RAM: %.2f%%", p.Pid, name, cpu, mem)

			// Выводим в консоль
			fmt.Println(msg)

			// Записываем в лог (с текущим временем)
			log.Printf("%s | %s\n", time.Now().Format(time.RFC3339), msg)

			// Отправляем сообщение через Telegram-бота конкретному пользователю
			SendAlertTo(chatID, p.Pid, name, cpu, mem)
		}
	}
}
