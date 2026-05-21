package  transport

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/services"
	"net/http"
	"errors"
	"strconv"
	"github.com/gin-gonic/gin"
)

type PromocodeHandler struct {
	promocode services.PromocodeService
}

func NewPromocodeHandler (promocode services.PromocodeService) *PromocodeHandler {
	return &PromocodeHandler{promocode: promocode}
}

func (h *PromocodeHandler) RegisterRoutes(router *gin.Engine) {
	promocodes := router.Group("/promocodes")
	{
		promocodes.POST("/", h.CreatePromo)
		promocodes.GET("/", h.GetAllPromo)
		promocodes.PATCH("/:id", h.UpdatePromo)
		promocodes.DELETE("/:id", h.DeletePromo)
		promocodes.POST("/validate", h.ValidatePromo)
	}
}

func (h *PromocodeHandler) CreatePromo(c *gin.Context) {
	var req models.PromocodeCreateRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promocode, err := h.promocode.CreatePromocode(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": promocode})
}

func (h *PromocodeHandler) GetAllPromo(c *gin.Context) {
	promocodes, err := h.promocode.GetAllPromocodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return	
	}

	c.JSON(http.StatusOK, gin.H{"data": promocodes})
}

func (h *PromocodeHandler) UpdatePromo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input models.PromocodeUpdateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promocode, err := h.promocode.UpdatePromocode(uint(id), input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPromocodeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return		
	}

	c.JSON(http.StatusOK, gin.H{"data": promocode})
}

func (h *PromocodeHandler) DeletePromo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.promocode.DeletePromocode(uint(id)); err != nil {
		switch {
		case errors.Is(err, services.ErrPromocodeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return		
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *PromocodeHandler) ValidatePromo(c *gin.Context) {
	var req models.PromocodeValidateRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promocode, err := h.promocode.ValidatePromocode(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPromocodeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPromocodeExpired), errors.Is(err, services.ErrPromocodeInactive):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return		
	}

	c.JSON(http.StatusOK, gin.H{"data": promocode})
}