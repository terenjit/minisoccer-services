package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "payment-service/common/error"
	"payment-service/constants"
	errConstant "payment-service/constants/error"
	errPayment "payment-service/constants/error/payment"
	"payment-service/domain/dto"
	"payment-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type IPaymentRepository interface {
	FindAllWithPagination(context.Context, *dto.PaymentRequestParam) ([]models.Payment, int64, error)
	FindByUUID(context.Context, string) (*models.Payment, error)
	FindByOrderID(context.Context, string) (*models.Payment, error)
	Create(context.Context, *gorm.DB, *dto.PaymentRequest) (*models.Payment, error)
	Update(context.Context, *gorm.DB, string, *dto.UpdatePaymentRequest) (*models.Payment, error)
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (p *PaymentRepository) FindAllWithPagination(c context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error) {
	var (
		payments []models.Payment
		sort     string
		total    int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *&param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit

	err := p.db.WithContext(c).Limit(limit).Offset(offset).Order(sort).Find(&payments).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = p.db.WithContext(c).Model(&payments).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return payments, total, nil

}

func (p *PaymentRepository) FindByUUID(c context.Context, uuid string) (*models.Payment, error) {
	var (
		payment models.Payment
	)

	err := p.db.WithContext(c).Where("uuid = ?", uuid).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &payment, nil

}

func (p *PaymentRepository) FindByOrderID(c context.Context, orderID string) (*models.Payment, error) {
	var (
		payment models.Payment
	)

	err := p.db.WithContext(c).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &payment, nil

}

func (p *PaymentRepository) Create(c context.Context, tx *gorm.DB, req *dto.PaymentRequest) (*models.Payment, error) {
	status := constants.Initial
	orderID := uuid.MustParse(req.OrderId)
	payment := models.Payment{
		UUID:        uuid.New(),
		OrderID:     orderID,
		Amount:      req.Amount,
		PaymentLink: req.PaymentLink,
		ExpiredAt:   req.ExpiredAt,
		Description: req.Description,
		Status:      &status,
	}

	err := tx.WithContext(c).Create(&payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &payment, nil

}

func (p *PaymentRepository) Update(c context.Context, tx *gorm.DB, orderId string, req *dto.UpdatePaymentRequest) (*models.Payment, error) {
	payment := models.Payment{
		Status:        req.Status,
		TransactionID: req.TransactionId,
		InvoiceLink:   req.InvoiceLink,
		PaidAt:        req.PaidAt,
		VANumber:      req.VANumber,
		Bank:          req.Bank,
		Acquirer:      req.Acquirer,
	}

	err := tx.WithContext(c).Where("order_id = ?", orderId).Updates(&payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &payment, nil
}
