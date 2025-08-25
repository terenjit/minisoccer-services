package repositories

import (
	"context"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	_ "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"

	"gorm.io/gorm"
)

type OrdertHistoryRepository struct {
	db *gorm.DB
}

type IOrderHistoryRepository interface {
	Create(context.Context, *gorm.DB, *dto.OrderHistoryRequest) error
}

func NewOrderHistoryRepository(db *gorm.DB) IOrderHistoryRepository {
	return &OrdertHistoryRepository{db: db}
}

func (o *OrdertHistoryRepository) Create(c context.Context, tx *gorm.DB, param *dto.OrderHistoryRequest) error {

	orderHistory := &models.OrderHistory{
		OrderID: param.OrderID,
		Status:  param.Status,
	}

	err := tx.WithContext(c).Create(&orderHistory).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}
