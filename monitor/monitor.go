package monitor

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func CheckProcesses(cpuThreshold, memThreshold float64) {
	logFile, err := os.OpenFile("alerts.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия файла лога: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	procs, err := process.Processes()
	if err != nil {
		log.Fatalf("Ошибка при получении процессов: %v", err)
	}

	for _, p := range procs {
		name, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()

		// Список системных процессов, дополнить позже
		var whitelist = map[string]bool{
			"System":      true,
			"Idle":        true,
			"svchost.exe": true,
		}

		if whitelist[name] {
			continue
		}

		if cpu > cpuThreshold || float64(mem) > memThreshold {
			msg := fmt.Sprintf("[⚠] PID: %d | Имя: %s | CPU: %.2f%% | RAM: %.2f%%", p.Pid, name, cpu, mem)
			fmt.Println(msg)
			// вывод в консоль
			logger.Printf("%s | %s\n", time.Now().Format(time.RFC3339), msg)
			// вывод на бота
			SendAlert(p.Pid, name, cpu, mem)
		}
	}
}
