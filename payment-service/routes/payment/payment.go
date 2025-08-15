package routes

import (
	"payment-service/clients"
	"payment-service/constants"
	controllers "payment-service/controllers/http"
	"payment-service/middlewares"

	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IPaymentRoute interface {
	Run()
}

func NewPaymentRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IPaymentRoute {
	return &PaymentRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *PaymentRoute) Run() {
	group := f.group.Group("/payment")
	group.POST("/webhook", f.controller.GetPayment().Webhook)
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetPayment().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetPayment().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{constants.Customer}, f.client), f.controller.GetPayment().Create)
}
