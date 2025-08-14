package repositories

import (
	repositoriesP "payment-service/repositories/payment"
	repositoriesPH "payment-service/repositories/paymenthistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetPayment() repositoriesP.IPaymentRepository
	GetPaymentHistory() repositoriesPH.IPaymentHistoryRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{
		db: db,
	}
}

func (r *Registry) GetPayment() repositoriesP.IPaymentRepository {
	return repositoriesP.NewPaymentRepository(r.db)
}

func (r *Registry) GetPaymentHistory() repositoriesPH.IPaymentHistoryRepository {
	return repositoriesPH.NewPaymentHistoryRepository(r.db)
}
