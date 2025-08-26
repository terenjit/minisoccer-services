package repositories

import (
	repoOrder "order-service/repositories/order"
	repoOrderField "order-service/repositories/orderfield"
	repoOrderHistory "order-service/repositories/orderhistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetOrder() repoOrder.IOrderRepository
	GetOrderHistory() repoOrderHistory.IOrderHistoryRepository
	GetOrderField() repoOrderField.IOrderFieldRepository
	GetTx() *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetOrder() repoOrder.IOrderRepository {
	return repoOrder.NewOrderRepository(r.db)
}
func (r *Registry) GetOrderHistory() repoOrderHistory.IOrderHistoryRepository {
	return repoOrderHistory.NewOrderHistoryRepository(r.db)
}
func (r *Registry) GetOrderField() repoOrderField.IOrderFieldRepository {
	return repoOrderField.NewOrderFieldRepository(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
