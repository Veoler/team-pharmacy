package main

import (
	"log"

	"github.com/Veoler/team-pharmacy/internal/config"
	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"

	// "github.com/Veoler/team-pharmacy/internal/repository"
	// "github.com/Veoler/team-pharmacy/internal/services"
	"github.com/Veoler/team-pharmacy/internal/transport"
)

func main() {
	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.OrderItem{}, &models.Cart{}, &models.CartItem{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, userRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, userRepo)

	transport.RegisterRoutes(router, userService, cartService, orderService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
