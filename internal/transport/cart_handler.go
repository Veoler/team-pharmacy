package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cart services.CartService
	user services.UserService
}

func NewCartHandler(
	cart services.CartService,
	user services.UserService,
) *CartHandler {
	return &CartHandler{cart: cart, user: user}
}

func (h *CartHandler) RegisterRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.POST("/:id/cart/items", h.AddItem)
		users.PATCH("/:id/cart/items/:item_id", h.AddQuantity) // добавлено
		users.DELETE("/:id/cart/items/:item_id", h.DeleteItem)
		users.DELETE("/:id/cart", h.DeleteCart)
	}
}

func (h *CartHandler) AddItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.CartCreateUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = new(uint)
	*req.UserID = uint(id)

	if _, err := h.user.GetByID(uint(id)); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cart.AddItem(req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddQuantity(ctx *gin.Context) {
	userIDp, userError := strconv.ParseUint(ctx.Param("id"), 10, 32)
	itemIDp, itemError := strconv.ParseUint(ctx.Param("item_id"), 10, 32)

	if userError != nil || itemError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "такие id не приветсвуются"})
		return
	}
	userID := uint(userIDp)
	itemID := uint(itemIDp)
	var req models.UpdateItemQuantityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cart.AddQuantity(&userID, &itemID, *req.Quantity)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrCartItemNotFound) {
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

	ctx.JSON(http.StatusOK, cart)
}

func (h *CartHandler) DeleteItem(ctx *gin.Context) {
	userIDp, userError := strconv.ParseUint(ctx.Param("id"), 10, 32)
	itemIDp, itemError := strconv.ParseUint(ctx.Param("item_id"), 10, 32)

	if userError != nil || itemError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "такие id не приветсвуются"})
		return
	}

	if err := h.cart.DeleteItem(uint(itemIDp), uint(userIDp)); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrCartItemNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *CartHandler) DeleteCart(ctx *gin.Context) {
	userIDp, userError := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if userError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "такие id не приветсвуются"})
		return
	}

	if err := h.cart.DeleteCart(uint(userIDp)); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
