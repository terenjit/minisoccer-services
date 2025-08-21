package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/clients/config"
	"order-service/common/util"
	config2 "order-service/config"
	"order-service/constants"
	"order-service/domain/dto"
	"time"

	"github.com/google/uuid"
)

type PaymentClient struct {
	client config.IClientConfig
}

type IPaymentClient interface {
	GetPaymentUUID(context.Context, uuid.UUID) (*PaymentData, error)
	CreatePaymentLink(context.Context, *dto.PaymentRequest) (*PaymentData, error)
}

func NewPaymentClient(client config.IClientConfig) IPaymentClient {
	return &PaymentClient{client: client}
}

func (p *PaymentClient) GetPaymentUUID(c context.Context, uuid uuid.UUID) (*PaymentData, error) {
	unixTime := time.Now().Unix()
	generateApikey := fmt.Sprintf("%s:%s:%d", config2.Cfg.AppName, p.client.SignatureKey(), unixTime)
	apiKey := util.GenerateSHA256(generateApikey)
	token := c.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response PaymentResponse
	request := p.client.Client().Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XserviceName, config2.Cfg.AppName).
		Set(constants.XrequestAt, fmt.Sprintf("%s", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/payment/%s", p.client.BaseURL(), uuid))

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment response: %s", response.Message)
	}

	return &response.Data, nil
}

func (p *PaymentClient) CreatePaymentLink(c context.Context, req *dto.PaymentRequest) (*PaymentData, error) {
	unixTime := time.Now().Unix()
	generateApikey := fmt.Sprintf("%s:%s:%d", config2.Cfg.AppName, p.client.SignatureKey(), unixTime)
	apiKey := util.GenerateSHA256(generateApikey)
	token := c.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, bodyResp, errs := p.client.Client().Clone().
		Post(fmt.Sprintf("%s/api/v1/payment", p.client.BaseURL())).
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XserviceName, config2.Cfg.AppName).
		Set(constants.XrequestAt, fmt.Sprintf("%s", unixTime)).
		Send(string(body)).End()

	if len(errs) > 0 {
		return nil, err
	}

	var response PaymentResponse
	if resp.StatusCode != http.StatusCreated {
		err = json.Unmarshal([]byte(bodyResp), &response)
		if err != nil {
			return nil, err
		}
		PaymentError := fmt.Errorf("payment response: %s", response.Message)
		return nil, PaymentError
	}

	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
