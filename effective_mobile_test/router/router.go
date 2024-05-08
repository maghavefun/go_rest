package router

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maghavefun/effective_mobile_test/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine *gin.Engine
}

type IRouter interface {
	Serve()
}

func NewRouter() *Router {
	return &Router{engine: gin.Default()}
}

func (r *Router) Serve() {
	f, err := os.Create("backend_info.log")
	if err != nil {
		log.Fatal("Error creating log file for gin:", err)
	}

	r.engine.Use(gin.LoggerWithWriter(f))

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var carsController controller.ICarsController = controller.NewCarController()

	v1 := r.engine.Group("/v1")
	{
		v1.GET("/cars", carsController.GetCars)
		v1.DELETE("/cars/:id", carsController.Delete)
		v1.PUT("/cars/:id", carsController.Update)
		v1.POST("/cars", carsController.Create)
	}

	r.engine.Run()
}
