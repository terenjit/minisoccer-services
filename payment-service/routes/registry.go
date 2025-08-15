package routes

import (
	"payment-service/clients"
	controllers "payment-service/controllers/http"
	routesF "payment-service/routes/payment"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IRegistry {
	return &Registry{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (r *Registry) PaymentRoute() routesF.IPaymentRoute {
	return routesF.NewPaymentRoute(r.controller, r.group, r.client)
}

func (r *Registry) Serve() {
	r.PaymentRoute().Run()
}
