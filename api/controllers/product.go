package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jeffvo/event-sourcing/internal/dto"
	usecase "github.com/jeffvo/event-sourcing/internal/usecases"
)

type Handler struct {
	usecase usecase.ProductUsecaseInterface
}

func NewProductHandler(usecase *usecase.ProductUsecaseInterface) *Handler {
	return &Handler{
		usecase: *usecase,
	}
}

func (handler *Handler) CreateProduct(context *gin.Context) {
	var input = new(dto.NewProductStockDTO)
	var err = context.BindJSON(input)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	result, err := handler.usecase.NewProductInStock(context.Request.Context(), input)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, result)
}

func (handler *Handler) GetProductById(context *gin.Context) {
	id, err := uuid.FromString(context.Param("id"))

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return

	}

	product, err := handler.usecase.GetProductStockById(context.Request.Context(), id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, product)
}

func (handler *Handler) RemoveStockFromProductById(context *gin.Context) {
	var input = new(dto.UpdateProductStockDto)
	var err = context.BindJSON(input)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	input.ID, err = uuid.FromString(context.Param("id"))

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return

	}

	err = handler.usecase.AdjustStock(context.Request.Context(), input)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "message": err.Error()})
		return
	}
}
