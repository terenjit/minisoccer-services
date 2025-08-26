package services

import (
	"context"
	"fmt"
	"order-service/clients"
	clientField "order-service/clients/field"
	clientPayment "order-service/clients/payment"
	clientUser "order-service/clients/user"
	"order-service/common/util"
	"order-service/constants"
	errOrder "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"order-service/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
}

type IOrderService interface {
	GetAllWithPagination(context.Context, *dto.OrderRequestParam) (*util.PaginationResult, error)
	GetByUUID(context.Context, string) (*dto.OrderResponse, error)
	GetOrderByUserId(context.Context) ([]dto.OrderByUserIDResponse, error)
	Create(context.Context, *dto.OrderRequest) (*dto.OrderResponse, error)
	HandlePayment(context.Context, *dto.PaymentData) error
}

func NewOrderService(repository repositories.IRepositoryRegistry, client clients.IClientRegistry) IOrderService {
	return &OrderService{repository: repository, client: client}
}

func (o *OrderService) GetAllWithPagination(c context.Context, param *dto.OrderRequestParam) (*util.PaginationResult, error) {
	orders, total, err := o.repository.GetOrder().FindAllWithPagination(c, param)
	if err != nil {
		return nil, err
	}
	orderResult := make([]dto.OrderResponse, 0, len(orders))
	for _, v := range orders {
		user, err := o.client.GetUser().GetUserByUUID(c, v.UserID)
		if err != nil {
			return nil, err
		}
		orderResult = append(orderResult, dto.OrderResponse{
			UUID:      v.UUID,
			Code:      v.Code,
			UserName:  user.Username,
			Amount:    v.Amount,
			Status:    v.Status.GetStatusString(),
			OrderDate: v.Date,
			CreatedAt: *v.CreatedAt,
			UpdatedAt: *v.UpdatedAt,
		})
	}

	paginationParam := util.PaginationParam{
		Page:  param.Page,
		Limit: param.Limit,
		Count: total,
		Data:  orderResult,
	}

	response := util.GeneratePagination(paginationParam)
	return &response, nil
}

func (o *OrderService) GetByUUID(c context.Context, uuid string) (*dto.OrderResponse, error) {
	var (
		v    *models.Order
		user *clientUser.UserData
		err  error
	)

	v, err = o.repository.GetOrder().FindByUUID(c, uuid)
	if err != nil {
		return nil, err
	}

	user, err = o.client.GetUser().GetUserByUUID(c, v.UserID)
	if err != nil {
		return nil, err
	}

	resp := dto.OrderResponse{
		UUID:      v.UUID,
		Code:      v.Code,
		UserName:  user.Username,
		Amount:    v.Amount,
		Status:    v.Status.GetStatusString(),
		OrderDate: v.Date,
		CreatedAt: *v.CreatedAt,
		UpdatedAt: *v.UpdatedAt,
	}

	return &resp, nil
}

func (o *OrderService) GetOrderByUserId(c context.Context) ([]dto.OrderByUserIDResponse, error) {
	var (
		order []models.Order
		user  = c.Value(constants.User).(*clientUser.UserData)
		err   error
	)
	order, err = o.repository.GetOrder().FindByUserID(c, user.UUID.String())
	if err != nil {
		return nil, err
	}

	orderLists := make([]dto.OrderByUserIDResponse, 0, len(order))
	for _, item := range order {
		payment, err := o.client.GetPayment().GetPaymentUUID(c, item.PaymentID)
		if err != nil {
			return nil, err
		}

		orderLists = append(orderLists, dto.OrderByUserIDResponse{
			Code:        item.Code,
			Amount:      fmt.Sprintf("%s", util.RupiahFormat(&item.Amount)),
			Status:      item.Status.GetStatusString(),
			OrderDate:   item.Date.String(),
			PaymentLink: payment.PaymentLink,
			InvoiceLink: payment.InvoiceLink,
		})
	}

	return orderLists, nil
}

func (o *OrderService) Create(c context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error) {
	var (
		order               *models.Order
		txErr, err          error
		user                = c.Value(constants.User).(*clientUser.UserData)
		field               *clientField.FieldData
		paymentResponse     *clientPayment.PaymentData
		orderFieldSchedules = make([]models.OrderField, 0, len(req.FieldScheduleIDs))
		totalAmount         float64
	)

	for _, fieldID := range req.FieldScheduleIDs {
		uuidParsed := uuid.MustParse(fieldID)
		field, err = o.client.GetField().GetFieldByUUID(c, uuidParsed)
		if err != nil {
			return nil, err
		}

		totalAmount += field.PricePerHour
		if field.Status == constants.BookedFieldStatus.String() {
			return nil, errOrder.ErrAlreadyBooked
		}
	}

	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		order, txErr = o.repository.GetOrder().Create(c, tx, &models.Order{
			UserID: user.UUID,
			Amount: totalAmount,
			Date:   time.Now(),
			Status: constants.Pending,
			IsPaid: false,
		})
		if txErr != nil {
			return txErr
		}

		for _, fieldID := range req.FieldScheduleIDs {
			uuidParsed := uuid.MustParse(fieldID)
			orderFieldSchedules = append(orderFieldSchedules, models.OrderField{
				OrderID:         order.ID,
				FieldScheduleID: uuidParsed,
			})
		}

		txErr = o.repository.GetOrderField().Create(c, tx, orderFieldSchedules)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(c, tx, &dto.OrderHistoryRequest{
			Status:  constants.Pending.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(time.Hour * 1)
		description := fmt.Sprintf("Payment Rent %s", field.FieldName)
		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(c, &dto.PaymentRequest{
			OrderID:     order.UUID,
			ExpiredAt:   expiredAt,
			Amount:      totalAmount,
			Description: description,
			CustomerDetail: dto.CustomerDetail{
				Name:  user.Name,
				Email: user.Email,
				Phone: user.PhoneNumber,
			},
			ItemDetails: []dto.ItemDetails{
				{
					ID:       uuid.New(),
					Name:     description,
					Amount:   totalAmount,
					Quantity: 1,
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrder().Update(c, tx, &models.Order{
			PaymentID: paymentResponse.UUID}, order.UUID)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:        order.UUID,
		Code:        order.Code,
		UserName:    user.Name,
		Amount:      order.Amount,
		Status:      order.Status.GetStatusString(),
		OrderDate:   order.Date,
		PaymentLink: paymentResponse.PaymentLink,
		CreatedAt:   *order.CreatedAt,
		UpdatedAt:   *order.UpdatedAt,
	}

	return &response, nil
}

func (o *OrderService) mapPaymentStatusToOrder(req *dto.PaymentData) (constants.OrderStatus, *models.Order) {

	var (
		status constants.OrderStatus
		order  *models.Order
	)

	switch req.Status {
	case constants.SettlementPaymentStatus:
		status = constants.PaymentSuccess
		order = &models.Order{
			IsPaid:    true,
			PaymentID: req.PaymentID,
			PaidAt:    req.PaidAt,
			Status:    status,
		}
	case constants.ExpiredPaymentStatus:
		status = constants.Expired
		order = &models.Order{
			IsPaid:    false,
			PaymentID: req.PaymentID,
			PaidAt:    req.PaidAt,
			Status:    status,
		}
	case constants.PendingPaymentStatus:
		status = constants.PendingPayment
		order = &models.Order{
			IsPaid:    false,
			PaymentID: req.PaymentID,
			PaidAt:    req.PaidAt,
			Status:    status,
		}
	}
	return status, order
}

func (o *OrderService) HandlePayment(c context.Context, req *dto.PaymentData) error {
	var (
		err, txErr          error
		order               *models.Order
		orderFieldSchedules []models.OrderField
	)

	status, body := o.mapPaymentStatusToOrder(req)
	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		txErr = o.repository.GetOrder().Update(c, tx, body, req.OrderID)
		if txErr != nil {
			return txErr
		}

		order, txErr = o.repository.GetOrder().FindByUUID(c, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(c, tx, &dto.OrderHistoryRequest{
			Status:  status.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		if req.Status == constants.SettlementPaymentStatus {
			orderFieldSchedules, txErr = o.repository.GetOrderField().FindByOrderID(c, order.ID)
			if txErr != nil {
				return txErr
			}

			fieldScheduleIDs := make([]string, 0, len(orderFieldSchedules))
			for _, item := range orderFieldSchedules {
				fieldScheduleIDs = append(fieldScheduleIDs, item.FieldScheduleID.String())
			}

			txErr = o.client.GetField().UpdateStatus(&dto.UpdateFieldScheduleStatusRequest{
				FieldScheduleIDs: fieldScheduleIDs,
			})
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
