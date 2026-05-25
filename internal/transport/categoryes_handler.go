package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/services"
	"github.com/gin-gonic/gin"
)

type CategoryesHandler struct {
	categoryes services.CategoryesService
}

func NewCategoryesHandler(
	categoryes services.CategoryesService,
) *CategoryesHandler {
	return &CategoryesHandler{
		categoryes: categoryes,
	}
}

func (h *CategoryesHandler) RegisterRoutes(router *gin.Engine) {
	categories := router.Group("/categories")
	{
		categories.GET("/", h.Get)
		categories.POST("/", h.Create)
		categories.GET("/:id/subcategories", h.GetSubcategoryes)
		categories.POST("/:id/subcategories", h.CreateSubcategory)
	}
}

func (h *CategoryesHandler) Get(ctx *gin.Context) {
	categoryes, err := h.categoryes.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"categories": categoryes})
}

func (h *CategoryesHandler) Create(ctx *gin.Context) {
	var req models.CategoryCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryes.Create(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created": category})
}

func (h *CategoryesHandler) GetSubcategoryes(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subcategories, err := h.categoryes.GetAllSubcategoryes(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"subcategories": subcategories})
}

func (h *CategoryesHandler) CreateSubcategory(ctx *gin.Context) {
	var req models.SubcategoryCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.CategoryID = new(uint)
	*req.CategoryID = uint(id)

	subcategory, err := h.categoryes.CreateSubcategory(req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created": subcategory})
}
