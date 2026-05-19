package  transport

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/services"
	"net/http"
	"errors"
	"strconv"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	review services.ReviewService
}

func NewReviewHandler (review services.ReviewService) *ReviewHandler {
	return &ReviewHandler{review: review}
}

func (h *ReviewHandler) RegisterRoutes(router *gin.Engine) {

	router.POST("/medicines/:id/reviews", h.CreateRev)
	router.GET("/medicines/:id/reviews", h.GetRevsFromMed)
	router.PATCH("/reviews/:id", h.UpdateRev)
	router.DELETE("/reviews/:id", h.DeleteRev)

}

func (h *ReviewHandler) CreateRev(c *gin.Context) {
	medicineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id лекарства"})
		return
	}

	var req models.ReviewCreateRequest
 
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.MedicineID = uint(medicineID)

	review, err := h.review.CreateReview(&req)
	if err != nil {
		if errors.Is(err, services.ErrNotPurchased) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(http.StatusCreated, gin.H{"data": review})
}

func (h *ReviewHandler) GetRevsFromMed(c *gin.Context) {
	medicineID, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный medicine_id"})
		return
	}
 
	reviews, err := h.review.GetReviewsFromMedicine(uint(medicineID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReviewNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrMedicineNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
 
	c.JSON(http.StatusOK, gin.H{"data": reviews})
}

func (h *ReviewHandler) UpdateRev(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
 
	var req models.ReviewUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	review, err := h.review.UpdateReview(uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReviewNotFound): 
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
 
	c.JSON(http.StatusOK, gin.H{"data": review})
}

func (h *ReviewHandler) DeleteRev(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
 
	if err := h.review.DeleteReview(uint(id)); err != nil {
		switch {
		case errors.Is(err, services.ErrReviewNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
 
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}