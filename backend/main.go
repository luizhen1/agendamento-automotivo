package main

import (
	"log"

	"agendamento-automotivo/backend/internal/bootstrap"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatalf("inicializar aplicação: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("executar aplicação: %v", err)
	}
}
