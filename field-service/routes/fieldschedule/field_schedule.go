package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldScheduleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldScheduleRoute interface {
	Run()
}

func NewFieldScheduleRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IFieldScheduleRoute {
	return &FieldScheduleRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *FieldScheduleRoute) Run() {
	group := f.group.Group("/field/schedule").Use(middlewares.AuthenticateWithoutToken())
	group.GET("", f.controller.GetFieldSchedule().GetAllFieldIdAndDate)
	group.GET("", f.controller.GetFieldSchedule().UpdateStatus)
	group.Use(middlewares.Authenticate())
	group.GET("/:uuid", middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetFieldSchedule().GetByUUID)
	group.GET("/pagination", middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetFieldSchedule().GetAllWithPagination)
	group.POST("/create", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Create)
	group.POST("/one-month", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().GenerateScheduleForOneMonth)
	group.PUT("/:uuid", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Update)
	group.DELETE("/:uuid", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Delete)

}
