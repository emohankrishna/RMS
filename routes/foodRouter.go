package routes

import (
	"github.com/emohankrishna/RMS/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/foods", controllers.GetFoods())
	incommingRoutes.GET("/foods/:food_id", controllers.GetFood())
	incommingRoutes.POST("/foods", controllers.CreateFood())
	incommingRoutes.PATCH("/foods/:food_id", controllers.UpdateFood())
}
