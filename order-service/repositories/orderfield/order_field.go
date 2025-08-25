package repositories

import (
	"context"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	_ "order-service/constants/error/order"
	"order-service/domain/models"

	"gorm.io/gorm"
)

type OrdertHistoryRepository struct {
	db *gorm.DB
}

type IOrderFieldRepository interface {
	FindByOrderID(context.Context, uint) ([]models.OrderField, error)
	Create(context.Context, *gorm.DB, []models.OrderField) error
}

func NewOrderFieldRepository(db *gorm.DB) IOrderFieldRepository {
	return &OrdertHistoryRepository{db: db}
}

func (o *OrdertHistoryRepository) FindByOrderID(c context.Context, OrderID uint) ([]models.OrderField, error) {
	var orderFields []models.OrderField

	err := o.db.WithContext(c).Where("order_id = ?", OrderID).Find(&orderFields).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orderFields, nil
}

func (o *OrdertHistoryRepository) Create(c context.Context, tx *gorm.DB, req []models.OrderField) error {

	err := tx.WithContext(c).Create(&req).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}
