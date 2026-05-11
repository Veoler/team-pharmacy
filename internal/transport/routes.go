package transport

import (
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"
	// "github.com/Veoler/team-pharmacy/internal/services"
)

func RegisterRoutes(
	router *gin.Engine,
	user services.UserService,
	cart services.CartService,
	order services.OrderService,
) {
	userHandler := NewUserHandler(user, order)
	userHandler.RegisterRoutes(router)
	cartHandler := NewCartHandler(cart, user)
	cartHandler.RegisterRoutes(router)
	orderHandler := NewOrderHandler(order)
	orderHandler.RegisterRoutes(router)
}
