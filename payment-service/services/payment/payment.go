package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	clients "payment-service/clients/midtrans"
	"payment-service/common/gcs"
	"payment-service/common/util"
	"payment-service/config"
	"payment-service/constants"
	errPayment "payment-service/constants/error/payment"
	"payment-service/controllers/kafka"
	"payment-service/domain/dto"
	"payment-service/domain/models"
	"payment-service/repositories"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PaymentService struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
	kafka      kafka.IKafkaRegistry
	midtrans   clients.IMidtransClient
}

type IPaymentService interface {
	GetAllWithPagination(context.Context, *dto.PaymentRequestParam) (*util.PaginationResult, error)
	GetByUUID(context.Context, string) (*dto.PaymentResponse, error)
	Create(context.Context, *dto.PaymentRequest) (*dto.PaymentResponse, error)
	Webhook(context.Context, *dto.Webhook) error
}

func NewPaymentService(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient, kafka kafka.IKafkaRegistry, midtrans clients.IMidtransClient) IPaymentService {
	return &PaymentService{
		repository: repository,
		gcs:        gcs,
		kafka:      kafka,
		midtrans:   midtrans,
	}
}

func (p *PaymentService) GetAllWithPagination(c context.Context, param *dto.PaymentRequestParam) (*util.PaginationResult, error) {
	payments, total, err := p.repository.GetPayment().FindAllWithPagination(c, param)
	if err != nil {
		return nil, err
	}
	paymentResult := make([]dto.PaymentResponse, 0, len(payments))
	for _, v := range payments {
		paymentResult = append(paymentResult, dto.PaymentResponse{
			UUID:          v.UUID,
			TransactionId: v.TransactionID,
			OrderID:       v.OrderID,
			Amount:        v.Amount,
			Status:        v.Status.GetStatusString(),
			PaymentLink:   v.PaymentLink,
			InvoiceLink:   v.InvoiceLink,
			VANumber:      v.VANumber,
			Bank:          v.Bank,
			Description:   v.Description,
			CreatedAt:     &v.CreatedAt,
			UpdatedAt:     &v.UpdatedAt,
			ExpireddAt:    &v.ExpiredAt,
		})
	}

	paginationParam := util.PaginationParam{
		Page:  param.Page,
		Limit: param.Limit,
		Count: total,
		Data:  paymentResult,
	}

	response := util.GeneratePagination(paginationParam)
	return &response, nil
}

func (p *PaymentService) GetByUUID(c context.Context, uuid string) (*dto.PaymentResponse, error) {
	payment, err := p.repository.GetPayment().FindByUUID(c, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		UUID:          payment.UUID,
		TransactionId: payment.TransactionID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status.GetStatusString(),
		PaymentLink:   payment.PaymentLink,
		InvoiceLink:   payment.InvoiceLink,
		VANumber:      payment.VANumber,
		Bank:          payment.Bank,
		Description:   payment.Description,
		CreatedAt:     &payment.CreatedAt,
		UpdatedAt:     &payment.UpdatedAt,
		ExpireddAt:    &payment.ExpiredAt,
	}, nil
}

func (p *PaymentService) Create(c context.Context, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	var (
		txErr, err error
		payment    *models.Payment
		response   *dto.PaymentResponse
		midtrans   *clients.MidtransData
	)

	err = p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		if !req.ExpiredAt.After(time.Now()) {
			return errPayment.ErrExpireAtInvalid
		}
		midtrans, txErr = p.midtrans.CreatePaymentLink(req)
		if txErr != nil {
			return txErr
		}

		paymentRequest := dto.PaymentRequest{
			OrderId:     req.OrderId,
			Amount:      req.Amount,
			Description: req.Description,
			ExpiredAt:   req.ExpiredAt,
			PaymentLink: midtrans.RedirectURL,
		}
		payment, txErr = p.repository.GetPayment().Create(c, tx, &paymentRequest)
		if txErr != nil {
			return txErr
		}

		txErr = p.repository.GetPaymentHistory().Create(c, tx, &dto.PaymentHistoryRequest{
			PaymentID: payment.ID,
			Status:    payment.Status.GetStatusString(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	response = &dto.PaymentResponse{
		UUID:        payment.UUID,
		OrderID:     payment.OrderID,
		Amount:      payment.Amount,
		Status:      payment.Status.GetStatusString(),
		PaymentLink: payment.PaymentLink,
		Description: payment.Description,
	}

	return response, nil
}

func (p *PaymentService) ConvertToIndonesianMonth(english string) string {
	monthMap := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"August":    "Agustus",
		"September": "Septembet",
		"October":   "Oktober",
		"December":  "Desember",
	}

	indonesianMonth, ok := monthMap[english]
	if !ok {
		return errors.New("month not found!").Error()
	}

	return indonesianMonth
}

func (p *PaymentService) generatePDF(req *dto.InvoiceRequest) ([]byte, error) {
	htmlTempPath := "template/invoice.html"
	htmpTemp, err := os.ReadFile(htmlTempPath)
	if err != nil {
		return nil, err
	}

	pdf, err := util.GeneratePDFfromHTML(string(htmpTemp))
	if err != nil {
		return nil, err
	}

	return pdf, nil
}

func (p *PaymentService) uploadToGCS(c context.Context, invoice string, pdf []byte) (string, error) {
	invoiceNumReplace := strings.ToLower(strings.ReplaceAll(invoice, "/", "-"))
	fileName := fmt.Sprintf("%s.pdf", invoiceNumReplace)
	url, err := p.gcs.UploadFile(c, fileName, pdf)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (p *PaymentService) RandomNum() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := random.Intn(900000) + 100000
	return number
}

func (p *PaymentService) mapTransactionStatusToEvent(status constants.PaymentStatusString) string {
	var paymentStatus string
	switch status {
	case constants.PendingString:
		paymentStatus = strings.ToUpper(constants.Pending.String())
	case constants.SettlementString:
		paymentStatus = strings.ToUpper(constants.Settlement.String())
	case constants.ExpireString:
		paymentStatus = strings.ToUpper(constants.Expire.String())
	case constants.InitialString:
		paymentStatus = strings.ToUpper(constants.Initial.String())
	}
	return paymentStatus
}

func (p *PaymentService) ProduceToKafka(req *dto.Webhook, payment *models.Payment, paidAt *time.Time) error {
	event := dto.KafkaEvent{
		Name: p.mapTransactionStatusToEvent(req.TransactionStatus),
	}

	metadata := dto.KafkaMetaData{
		Sender:    "payment-service",
		SendingAt: time.Now().Format(time.RFC3339),
	}

	body := dto.KafkaBody{
		Type: "JSON",
		Data: &dto.KafkaData{
			OrderID:   payment.OrderID,
			PaymentID: payment.UUID,
			Status:    string(req.TransactionStatus),
			PaidAt:    paidAt,
			ExpiredAt: payment.ExpiredAt,
		},
	}

	kafkaMsg := dto.KafkaMessage{
		Event:    event,
		Metadata: metadata,
		Body:     body,
	}

	topic := config.Cfg.Kafka.Topic
	kafkaMsgJSON, _ := json.Marshal(kafkaMsg)
	err := p.kafka.GetKafkaProducer().ProduceMessage(topic, kafkaMsgJSON)
	if err != nil {
		return err
	}
	return nil
}

func (p *PaymentService) Webhook(c context.Context, req *dto.Webhook) error {
	var (
		txErr, err         error
		paymentAfterUpdate *models.Payment
		paidAt             *time.Time
		invoiceLink        string
		pdf                []byte
	)

	err = p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		_, txErr = p.repository.GetPayment().FindByOrderID(c, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		if req.TransactionStatus == constants.SettlementString {
			now := time.Now()
			paidAt = &now
		}

		status := req.TransactionStatus.GetStatusInt()
		vaNumber := req.VANumbers[0].VaNumber
		bank := req.VANumbers[0].Bank
		_, txErr = p.repository.GetPayment().Update(c, tx, req.OrderID.String(), &dto.UpdatePaymentRequest{
			TransactionId: &req.TransactionID,
			Status:        &status,
			PaidAt:        paidAt,
			VANumber:      &vaNumber,
			Bank:          &bank,
			Acquirer:      req.Acquirer,
		})
		if txErr != nil {
			return txErr
		}

		paymentAfterUpdate, txErr = p.repository.GetPayment().FindByOrderID(c, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		txErr = p.repository.GetPaymentHistory().Create(c, tx, &dto.PaymentHistoryRequest{
			PaymentID: paymentAfterUpdate.ID,
			Status:    paymentAfterUpdate.Status.GetStatusString(),
		})

		if req.TransactionStatus == constants.SettlementString {
			paidDay := paidAt.Format("02")
			paidMonth := p.ConvertToIndonesianMonth(paidAt.Format("January"))
			paidYear := paidAt.Format("2006")
			invoiceNum := fmt.Sprintf("INV/%s/ORD/%d", time.Now().Format(time.DateOnly), p.RandomNum())
			total := util.RupiahFormat(&paymentAfterUpdate.Amount)
			invoiceReq := dto.InvoiceRequest{
				InvoiceNumber: invoiceNum,
				Data: dto.InvoiceData{
					PaymentDetail: dto.InvoicePaymentDetail{
						PaymentMethod: req.PaymentType,
						BankName:      strings.ToUpper(*paymentAfterUpdate.Bank),
						VANumber:      *paymentAfterUpdate.VANumber,
						Date:          fmt.Sprintf("%s %s %s", paidDay, paidMonth, paidYear),
						IsPaid:        true,
					},
					Items: []dto.InvoiceItem{
						{
							Description: *paymentAfterUpdate.Description,
							Price:       total,
						},
					},
					Total: total,
				},
			}
			pdf, txErr = p.generatePDF(&invoiceReq)
			if txErr != nil {
				return txErr
			}

			invoiceLink, txErr = p.uploadToGCS(c, invoiceNum, pdf)
			if txErr != nil {
				return txErr
			}

			paymentAfterUpdate, txErr = p.repository.GetPayment().Update(c, tx, req.OrderID.String(), &dto.UpdatePaymentRequest{
				InvoiceLink: &invoiceLink,
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

	paymentAfterUpdate, txErr = p.repository.GetPayment().FindByOrderID(c, req.OrderID.String())
	if txErr != nil {
		return txErr
	}

	err = p.ProduceToKafka(req, paymentAfterUpdate, paidAt)
	if err != nil {
		return err
	}

	return nil
}
