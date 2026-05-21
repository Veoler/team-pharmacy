package transport

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/services"
	"net/http"
	"errors"
	"strconv"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	payment services.PaymentService
}

func NewPaymentHandler(payment services.PaymentService) *PaymentHandler {
	return &PaymentHandler{payment: payment}
}

func (h *PaymentHandler) RegisterRoutes (router *gin.Engine) {

	router.POST("/orders/:id/payments", h.CreatePay)
	router.GET("/orders/:id/payments", h.GetPayFromOrd)
	router.GET("/payments/:id", h.GetPayByID)
	router.DELETE("/payments/:id", h.DeletePay)	

}

func (h *PaymentHandler) CreatePay(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.PaymentCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OrderID = uint(orderID)
	
	payment, summary, err := h.payment.CreatePayment(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPaymentExceedsTotal):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"payment": payment, "order": summary})
}

func (h *PaymentHandler) GetPayFromOrd (c *gin.Context) {
	orderID,  err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payments, err := h.payment.GetPaymentFromOrder(uint(orderID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrOrderNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPaymentNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}

func (h *PaymentHandler) GetPayByID (c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.payment.GetPaymentByID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPaymentNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func (h *PaymentHandler) DeletePay (c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.payment.DeletePayment(uint(id)); err != nil {
		switch {
        case errors.Is(err, services.ErrPaymentNotFound):
        	    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    	default:
        	    c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
        }
        return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}


