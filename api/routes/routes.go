package routes

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gin-gonic/gin"

	controller "github.com/jeffvo/event-sourcing/api/controllers"
	repository "github.com/jeffvo/event-sourcing/internal/repositories"
	usecase "github.com/jeffvo/event-sourcing/internal/usecases"
)

// Function that is used to register all the endpoints that are used by the application
func RegisterHTTPEndpoints(db *esdb.Client) *gin.Engine {
	var router = gin.Default()
	var RouterGroup = router.Group("")

	var repository = repository.CreateProductRepository(db)
	var usecase = usecase.CreateProductUsecase(repository)
	var handler = controller.NewProductHandler(&usecase)

	RouterGroup.GET("/products/:id", handler.GetProductById)
	RouterGroup.POST("/products/", handler.CreateProduct)
	RouterGroup.PUT("/products/:id", handler.RemoveStockFromProductById)

	return router
}
