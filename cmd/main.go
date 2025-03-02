package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "Yandex_Calc_V2.0/docs" // which is the generated folder after swag init

	"Yandex_Calc_V2.0/internal/app"
)

// cmd/main.go
//
// @title Yandex Calculator API
// @version 1.0
// @description Calculator service with distributed computation
// @termsOfService http://swagger.io/terms/
//
// @contact.name Nyamerka
//
// @host localhost:8080
// @BasePath /api/v1
func main() {
	orchestrator := app.NewOrchestrator()

	go func() {
		if err := orchestrator.StartServer(); err != nil {
			log.Fatalf("Failed to start orchestrator: %v", err)
		}
	}()

	agent := app.NewAgent()

	go agent.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutting down orchestrator and agent...")

	time.Sleep(2 * time.Second)
	log.Println("Orchestrator and agent shutdown complete.")
}
