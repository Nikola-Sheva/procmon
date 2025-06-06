package main

import (
	"flag"
	"procmon/monitor"
	"time"
)

func main() {
	cpuThreshold := flag.Float64("cpu", 30.0, "Порог CPU (%)")
	memThreshold := flag.Float64("mem", 10.0, "Порог памяти (%)")
	interval := flag.Int("interval", 10, "Интервал мониторинга (сек)")

	monitor.InitBot("7941396107:AAESmpNfy5YVDQ7zelunO-P9UOCBTHpmrDw")

	for {
		monitor.CheckProcesses(*cpuThreshold, *memThreshold)
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
