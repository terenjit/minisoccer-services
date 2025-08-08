package services

import (
	"field-service/common/gcs"
	"field-service/repositories"
	servicesField "field-service/services/field"
	servicesFieldSchedule "field-service/services/fieldschedule"
	servicesTime "field-service/services/time"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

type IServiceRegistry interface {
	GetField() servicesField.IfieldService
	GetFieldSchedule() servicesFieldSchedule.IFieldScheduleService
	GetTime() servicesTime.ITimeService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IServiceRegistry {
	return &Registry{repository: repository, gcs: gcs}
}

func (r *Registry) GetField() servicesField.IfieldService {
	return servicesField.NewFieldService(r.repository, r.gcs)
}

func (r *Registry) GetFieldSchedule() servicesFieldSchedule.IFieldScheduleService {
	return servicesFieldSchedule.NewFieldScheduleService(r.repository)
}

func (r *Registry) GetTime() servicesTime.ITimeService {
	return servicesTime.NewTimeService(r.repository)
}
