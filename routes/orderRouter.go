package routes

import (
	"github.com/emohankrishna/RMS/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/orders", controllers.GetOrders())
	incommingRoutes.GET("/orders/:order_id", controllers.GetOrder())
	incommingRoutes.POST("/orders", controllers.CreateOrder())
	incommingRoutes.PATCH("/orders/:order_id", controllers.UpdateOrder())
}
