package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	order services.OrderService
}

func NewOrderHandler(
	order services.OrderService,
) *OrderHandler {
	return &OrderHandler{order: order}
}

func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/users/:id/orders", h.CreateOrder)
	router.GET("/orders/:id", h.GetByID)
	router.PATCH("/orders/:id/status", h.UpdateStatus)
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.OrderCreateRequest
////////////////////////////////////////////////////////////
	// if err := ctx.ShouldBindJSON(id); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	req.UserID = new(uint)
	*req.UserID = uint(id)

	order, err := h.order.CreateOrder(req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrCartNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created": order})
}

func (h *OrderHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.order.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req models.OrderUpdateStatusRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.order.UpdateOrderStatus(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status updated": order})
}
