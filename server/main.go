package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	CPUUsage    float64   `json:"cpu"`
	MemoryUsage float64   `json:"memory"`
	DiskUsage   float64   `json:"disk"`
	Timestamp   time.Time `json:"timestamp"`
}

var latestMetrics Metrics
var mutex sync.RWMutex

func main() {
	go startTCPServer()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/metrics", metricsHandler)

	fmt.Println("HTTP Server Started")
	fmt.Println("http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

func startTCPServer() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	fmt.Println("TCP Server Started")
	fmt.Println("Listening on localhost:8081")

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	decoder := json.NewDecoder(connection)
	for {
		var metrics Metrics
		err := decoder.Decode(&metrics)
		if err != nil {
			return
		}
		mutex.Lock()
		latestMetrics = metrics
		mutex.Unlock()
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.RLock()
	defer mutex.RUnlock()
	json.NewEncoder(w).Encode(latestMetrics)
}
