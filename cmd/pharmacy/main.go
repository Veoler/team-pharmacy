package main

import (
	"log"

	"github.com/Veoler/team-pharmacy/internal/config"
	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/Veoler/team-pharmacy/internal/transport"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.OrderItem{}, &models.Cart{}, &models.CartItem{}, &models.Payment{}, &models.Promocode{}, &models.Review{}, &models.Category{}, &models.Subcategory{}, &models.Medicine{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	promocodeRepo := repository.NewPromocodeRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	categoryRepo := repository.NewCategoryesRepository(db)
	medicineRepo := repository.NewMedicineRepository(db)

	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, userRepo, medicineRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, userRepo, promocodeRepo, medicineRepo)
	paymentService := services.NewPaymentService(paymentRepo, orderRepo)
	promocodeService := services.NewPromocodeService(promocodeRepo)
	reviewService := services.NewReviewService(reviewRepo, orderRepo, medicineRepo)
	categoryService := services.NewCategoryesService(categoryRepo)
	medicineService := services.NewMedicineService(medicineRepo, categoryRepo)

	transport.RegisterRoutes(router, userService, cartService, orderService, paymentService, promocodeService, reviewService, categoryService, medicineService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
