package main

import (
	"log"

	"github.com/Veoler/team-pharmacy/internal/config"
	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/Veoler/team-pharmacy/internal/transport"
)

func main() {
	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.OrderItem{}, &models.Cart{}, &models.CartItem{}, &models.Payment{}, &models.Promocode{}, &models.Review{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	promocodeRepo := repository.NewPromocodeRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, userRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, userRepo)
	paymentService := services.NewPaymentService(paymentRepo, orderRepo)
	promocodeService := services.NewPromocodeService(promocodeRepo)
	reviewService := services.NewReviewService(reviewRepo, orderRepo, medicineRepo)

	transport.RegisterRoutes(router, userService, cartService, orderService, paymentService, promocodeService, reviewService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
