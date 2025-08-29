package routes

import (
	"order-service/clients"
	controllers "order-service/controllers/http"
	routes "order-service/routes/order"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
	group      *gin.RouterGroup
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, client clients.IClientRegistry, group *gin.RouterGroup) IRouteRegistry {
	return &Registry{controller: controller, client: client, group: group}
}

func (r *Registry) Serve() {
	r.OrderRoute().Run()
}

func (r *Registry) OrderRoute() routes.IOrderRoute {
	return routes.NewOrderRoute(r.controller, r.client, r.group)
}
