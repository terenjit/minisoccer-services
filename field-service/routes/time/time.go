package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type TimeRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type ITimeRoute interface {
	Run()
}

func NewTimeRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) ITimeRoute {
	return &TimeRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *TimeRoute) Run() {
	group := f.group.Group("/time").Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetTime().GetAll)
	group.GET("/:uuid", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetTime().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetTime().Create)

}
