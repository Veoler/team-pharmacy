package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"
)

type MedicineHandler struct {
	medicine services.MedicineService
}

func NewMedicineHandler(
	medicine services.MedicineService,
) *MedicineHandler {
	return &MedicineHandler{
		medicine: medicine,
	}
}

func (h *MedicineHandler) RegisterRoutes(router *gin.Engine) {
	medicines := router.Group("/medicines")
	{
		medicines.GET("/", h.Get)
		medicines.GET("/:id", h.GetByID)
		medicines.POST("/", h.Create)
		medicines.PATCH("/:id", h.Update)
		medicines.DELETE("/:id", h.Delete)
	}
}

func (h *MedicineHandler) Get(ctx *gin.Context) {
	medicines, err := h.medicine.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"medicines": medicines})
}

func (h *MedicineHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	medicine, err := h.medicine.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrMedicineNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, medicine)
}

func (h *MedicineHandler) Create(ctx *gin.Context) {
	var req models.MedicineCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medicine, err := h.medicine.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSubCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created": medicine})
}

func (h *MedicineHandler) Update(ctx *gin.Context) {
	var req models.MedicineUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medicine, err := h.medicine.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSubCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"edited": medicine})
}

func (h *MedicineHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.medicine.Delete(uint(id)); err != nil {
		if errors.Is(err, services.ErrMedicineNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
