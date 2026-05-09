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
) {
	userHandler := NewUserHandler(user)
	userHandler.RegisterRoutes(router)
	cartHandler := NewCartHandler(cart, user)
	cartHandler.RegisterRoutes(router)
}
