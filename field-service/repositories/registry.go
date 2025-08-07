package repositories

import (
	repoField "field-service/repositories/field"
	repoFieldSchedule "field-service/repositories/fieldschedule"
	repoTime "field-service/repositories/time"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetField() repoField.IFieldRepository
	GetFieldSchedule() repoFieldSchedule.IFieldScheduleRepository
	GetTime() repoTime.ITimeRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetField() repoField.IFieldRepository {
	return repoField.NewFieldRepository(r.db)
}
func (r *Registry) GetFieldSchedule() repoFieldSchedule.IFieldScheduleRepository {
	return repoFieldSchedule.NewFieldScheduleRepository(r.db)
}
func (r *Registry) GetTime() repoTime.ITimeRepository {
	return repoTime.NewTimeRepository(r.db)
}
