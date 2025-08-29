package routes

import (
	"order-service/clients"
	"order-service/constants"
	controllers "order-service/controllers/http"
	"order-service/middlewares"

	"github.com/gin-gonic/gin"
)

type OrderRoute struct {
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
	group      *gin.RouterGroup
}

type IOrderRoute interface {
	Run()
}

func NewOrderRoute(controller controllers.IControllerRegistry, client clients.IClientRegistry, group *gin.RouterGroup) IOrderRoute {
	return &OrderRoute{controller: controller, client: client, group: group}
}

func (o *OrderRoute) Run() {
	group := o.group.Group("/order")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{constants.Admin, constants.Customer}, o.client), o.controller.GetOrder().GetAllWIthPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{constants.Admin, constants.Customer}, o.client), o.controller.GetOrder().GetByUUID)
	group.GET("/user", middlewares.CheckRole([]string{constants.Customer}, o.client), o.controller.GetOrder().GetOrderByUserID)
	group.POST("", middlewares.CheckRole([]string{constants.Customer}, o.client), o.controller.GetOrder().Create)
}
