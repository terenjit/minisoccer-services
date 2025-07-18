package services

import (
	"user-service/repositories"
	services "user-service/services/user"
)

type Registry struct {
	repository repositories.IRegistryRepository
}

type IServiceRegistry interface {
	GetUser() services.IUserService
}

func NewServiceRegistry(repository repositories.IRegistryRepository) IServiceRegistry {
	return &Registry{repository: repository}
}

func (r *Registry) GetUser() services.IUserService {
	return services.NewUserService(r.repository)
}
