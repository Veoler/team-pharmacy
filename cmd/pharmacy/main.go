package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Veoler/team-pharmacy/internal/config"
	"github.com/Veoler/team-pharmacy/internal/models"
	// "github.com/Veoler/team-pharmacy/internal/repository"
	// "github.com/Veoler/team-pharmacy/internal/services"
	"github.com/Veoler/team-pharmacy/internal/transport"
)

func main() {
	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.Payment{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	router := gin.Default()

	transport.RegisterRoutes(router)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
