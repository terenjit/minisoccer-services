package clients

import (
	errPayment "payment-service/constants/error/payment"
	"payment-service/domain/dto"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type MidtransClient struct {
	ServerKey    string
	IsProduction bool
}

type IMidtransClient interface {
	CreatePaymentLink(*dto.PaymentRequest) (*MidtransData, error)
}

func NewMidtransClient(serverKey string, isProduction bool) *MidtransClient {
	return &MidtransClient{
		ServerKey:    serverKey,
		IsProduction: isProduction,
	}
}

func (c *MidtransClient) CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error) {
	var (
		snapClient   snap.Client
		isProduction = midtrans.Sandbox
	)

	expiryDateTime := request.ExpiredAt
	currentTime := time.Now()
	timeDiff := expiryDateTime.Sub(currentTime)
	if timeDiff <= 0 {
		logrus.Errorf("expired at invalid")
		return nil, errPayment.ErrExpireAtInvalid
	}

	expiryUnit := "minute"
	expiryDuration := int64(timeDiff.Minutes())

	if timeDiff.Hours() >= 1 {
		expiryUnit = "hour"
		expiryDuration = int64(timeDiff.Hours())
	} else if timeDiff.Hours() >= 24 {
		expiryUnit = "day"
		expiryDuration = int64(timeDiff.Hours() / 24)
	}

	if c.IsProduction {
		isProduction = midtrans.Production
	}

	snapClient.New(c.ServerKey, isProduction)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderId,
			GrossAmt: int64(request.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.CustomerDetail.Name,
			Email: request.CustomerDetail.Email,
			Phone: request.CustomerDetail.Phone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    request.ItemDetail[0].ID,
				Name:  request.ItemDetail[0].Name,
				Price: int64(request.ItemDetail[0].Amount),
				Qty:   int32(request.ItemDetail[0].Quantity),
			},
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     expiryUnit,
			Duration: expiryDuration,
		},
	}

	response, err := snapClient.CreateTransaction(req)
	if err != nil {
		logrus.Errorf("Error create transaction midtrans: %v", err)
		return nil, err
	}

	return &MidtransData{
		RedirectURL: response.RedirectURL,
		Token:       response.Token,
	}, nil

}
