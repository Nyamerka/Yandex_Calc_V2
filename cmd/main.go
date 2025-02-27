package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Yandex_Calc_V2.0/internal/app"
)

func main() {
	orchestrator := app.NewOrchestrator()

	go func() {
		if err := orchestrator.StartServer(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	agent := app.SetDefaultAgent()

	go agent.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutting down server and agent...")

	time.Sleep(2 * time.Second)
	log.Println("Server and agent shutdown complete.")
}
