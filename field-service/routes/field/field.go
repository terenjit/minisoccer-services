package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldRoute interface {
	Run()
}

func NewFieldRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IFieldRoute {
	return &FieldRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *FieldRoute) Run() {
	// Public routes (no token required)
	publicGroup := f.group.Group("/field").Use(middlewares.AuthenticateWithoutToken())
	publicGroup.GET("", f.controller.GetField().GetAllWithoutPagination)
	publicGroup.GET("/:uuid", f.controller.GetField().GetByUUID)

	// Protected routes (authentication + role check)
	protectedGroup := f.group.Group("/field").Use(middlewares.Authenticate())
	protectedGroup.GET("/pagination",
		middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client),
		f.controller.GetField().GetAllWithPagination,
	)
	protectedGroup.POST("/create",
		middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client),
		f.controller.GetField().Create,
	)
	protectedGroup.PUT("/:uuid",
		middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client),
		f.controller.GetField().Update,
	)
	protectedGroup.DELETE("/:uuid",
		middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client),
		f.controller.GetField().Delete,
	)
}
