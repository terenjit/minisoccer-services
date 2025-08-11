package routes

import (
	"field-service/clients"
	"field-service/controllers"
	routesF "field-service/routes/field"
	routesFS "field-service/routes/fieldschedule"
	routesT "field-service/routes/time"

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

func (r *Registry) fieldRoute() routesF.IFieldRoute {
	return routesF.NewFieldRoute(r.controller, r.group, r.client)
}

func (r *Registry) fieldScheduleRoute() routesFS.IFieldScheduleRoute {
	return routesFS.NewFieldScheduleRoute(r.controller, r.group, r.client)
}

func (r *Registry) timeRoute() routesT.ITimeRoute {
	return routesT.NewTimeRoute(r.controller, r.group, r.client)
}

func (r *Registry) Serve() {
	r.fieldRoute().Run()
	r.fieldScheduleRoute().Run()
	r.timeRoute().Run()
}
