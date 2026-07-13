package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type Metrics struct {
	CPUUsage    float64   `json:"cpu"`
	MemoryUsage float64   `json:"memory"`
	DiskUsage   float64   `json:"disk"`
	Timestamp   time.Time `json:"timestamp"`
}

func collectMetrics() (Metrics, error) {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return Metrics{}, err
	}

	memoryUsage, err := mem.VirtualMemory()
	if err != nil {
		return Metrics{}, err
	}

	diskUsage, err := disk.Usage("/")
	if err != nil {
		return Metrics{}, err
	}

	return Metrics{
		CPUUsage:    cpuUsage[0],
		MemoryUsage: memoryUsage.UsedPercent,
		DiskUsage:   diskUsage.UsedPercent,
		Timestamp:   time.Now(),
	}, nil
}

func main() {
	connection, err := net.Dial("tcp", "192.168.18.6:8081")
	if err != nil {
		fmt.Println("Unable to connect to server:", err)
		return
	}
	defer connection.Close()

	encoder := json.NewEncoder(connection)

	for {
		metrics, err := collectMetrics()
		if err != nil {
			fmt.Println("Error collecting metrics:", err)
		}

		err = encoder.Encode(metrics)
		if err != nil {
			fmt.Println("Error sending metrics:", err)
			return
		}

		time.Sleep(5 * time.Second) // sends metrics every 5 seconds
	}
}
