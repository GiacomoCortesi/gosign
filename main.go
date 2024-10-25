package main

import (
	"log"

	"github.com/GiacomoCortesi/gosign/api"
	"github.com/GiacomoCortesi/gosign/persistence"
	"github.com/GiacomoCortesi/gosign/service"
)

const (
	ListenAddress = ":8080"
)

func main() {
	service := service.NewSignatureDeviceService(persistence.NewInMemorySignatureDeviceRepository())
	server := api.NewServer(ListenAddress, service)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
