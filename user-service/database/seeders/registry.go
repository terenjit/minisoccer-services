package seeders

import "gorm.io/gorm"

type Registry struct {
	db *gorm.DB
}

type ISseederRegistry interface {
	Run()
}

func NewSeederRegistry(db *gorm.DB) ISseederRegistry {
	return &Registry{db: db}
}

func (s *Registry) Run() {
	RunRoleSeeder(s.db)
	RunUserSeeder(s.db)
}
