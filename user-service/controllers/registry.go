package controllers

import (
	controllers "user-service/controllers/user"
	"user-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IUserControllerRegistry interface {
	GetUserController() controllers.IUserController
}

func NewControllerREgistry(service services.IServiceRegistry) IUserControllerRegistry {
	return &Registry{service: service}
}

func (h *Registry) GetUserController() controllers.IUserController {
	return controllers.NewUserController(h.service)
}
