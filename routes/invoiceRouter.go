package routes

import (
	"github.com/emohankrishna/RMS/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/invoices", controllers.GetInvoices())
	incommingRoutes.GET("/invoices/:invoive_id", controllers.GetInvoice())
	incommingRoutes.POST("/invoices", controllers.CreateInvoice())
	incommingRoutes.PATCH("/invoices/:invoice_id", controllers.UpdateInvoice())
}
