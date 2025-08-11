package controllers

import (
	controllersF "field-service/controllers/field"
	controllersFS "field-service/controllers/fieldschedule"
	controllersT "field-service/controllers/time"
	"field-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetField() controllersF.IFieldController
	GetFieldSchedule() controllersFS.IFieldScheduleController
	GetTime() controllersT.ITimeController
}

func NewServiceRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (r *Registry) GetField() controllersF.IFieldController {
	return controllersF.NewFieldController(r.service)
}

func (r *Registry) GetFieldSchedule() controllersFS.IFieldScheduleController {
	return controllersFS.NewFieldScheduleController(r.service)
}

func (r *Registry) GetTime() controllersT.ITimeController {
	return controllersT.NewTimeController(r.service)
}
