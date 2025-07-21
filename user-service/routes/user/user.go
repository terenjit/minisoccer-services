package routes

import (
	"user-service/controllers"
	"user-service/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.IUserControllerRegistry
	group      *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserROute(controller controllers.IUserControllerRegistry, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{controller: controller, group: group}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/auth")
	group.GET("/user", middlewares.Authenticate(), u.controller.GetUserController().GetUserLogin)
	group.GET("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().GetUserByUUID)
	group.POST("/login", u.controller.GetUserController().Login)
	group.POST("/register", u.controller.GetUserController().Register)
	group.PUT("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().Update)
}
