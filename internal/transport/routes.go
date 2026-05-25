package transport

import (
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	user services.UserService,
	cart services.CartService,
	order services.OrderService,
	payment services.PaymentService,
	promocode services.PromocodeService,
	review services.ReviewService,
	category services.CategoryesService,
	mediicne services.MedicineService,

) {
	userHandler := NewUserHandler(user, order)
	userHandler.RegisterRoutes(router)

	cartHandler := NewCartHandler(cart, user)
	cartHandler.RegisterRoutes(router)

	orderHandler := NewOrderHandler(order)
	orderHandler.RegisterRoutes(router)

	paymentHandler := NewPaymentHandler(payment)
	paymentHandler.RegisterRoutes(router)

	promocodeHandler := NewPromocodeHandler(promocode)
	promocodeHandler.RegisterRoutes(router)

	reviewHandler := NewReviewHandler(review)
	reviewHandler.RegisterRoutes(router)

	categoryHandler := NewCategoryesHandler(category)
	categoryHandler.RegisterRoutes(router)

	medicineHandler := NewMedicineHandler(mediicne)
	medicineHandler.RegisterRoutes(router)
}
