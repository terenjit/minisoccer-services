package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/clients/config"
	"order-service/common/util"
	Cfg "order-service/config"
	"order-service/constants"
	"order-service/domain/dto"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FieldClient struct {
	client config.IClientConfig
}

type IFieldClient interface {
	GetFieldByUUID(context.Context, uuid.UUID) (*FieldData, error)
	UpdateStatus(*dto.UpdateFieldScheduleStatusRequest) error
}

func NewFieldClient(client config.IClientConfig) IFieldClient {
	return &FieldClient{client: client}
}

func (f *FieldClient) GetFieldByUUID(c context.Context, uuid uuid.UUID) (*FieldData, error) {
	unixTime := time.Now().Unix()
	generateApikey := fmt.Sprintf("%s:%s:%d", "field-services", Cfg.Cfg.InternalService.Field.SignatureKey, unixTime)
	apiKey := util.GenerateSHA256(generateApikey)
	token := c.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/field/schedule/%s", Cfg.Cfg.InternalService.Field.Host, uuid), nil)
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.XrequestAt, fmt.Sprintf("%d", unixTime))
	req.Header.Set(constants.XserviceName, "field-services")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("err: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var response FieldResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return &response.Data, nil
}

// func (f *FieldClient) GetFieldByUUID(c context.Context, uuid uuid.UUID) (*FieldData, error) {
// 	unixTime := time.Now().Unix()
// 	generateApikey := fmt.Sprintf("%s:%s:%d", config2.Cfg.AppName, f.client.SignatureKey(), unixTime)
// 	apiKey := util.GenerateSHA256(generateApikey)
// 	token := c.Value(constants.Token).(string)
// 	bearerToken := fmt.Sprintf("Bearer %s", token)

// 	var response FieldResponse
// 	request := f.client.Client().Set(constants.Authorization, bearerToken).
// 		Set(constants.XApiKey, apiKey).
// 		Set(constants.XserviceName, config2.Cfg.AppName).
// 		Set(constants.XrequestAt, fmt.Sprintf("%s", unixTime)).
// 		Get(fmt.Sprintf("%s/api/v1/field/schedule/%s", f.client.BaseURL(), uuid))

// 	resp, _, errs := request.EndStruct(&response)
// 	if len(errs) > 0 {
// 		return nil, errs[0]
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("payment response: %s", response.Message)
// 	}

// 	return &response.Data, nil
// }

// func (f *FieldClient) UpdateStatus(req *dto.UpdateFieldScheduleStatusRequest) error {
// 	unixTime := time.Now().Unix()
// 	generateApikey := fmt.Sprintf("%s:%s:%d", config2.Cfg.AppName, f.client.SignatureKey(), unixTime)
// 	apiKey := util.GenerateSHA256(generateApikey)

// 	body, err := json.Marshal(req)
// 	if err != nil {
// 		return err
// 	}

// 	resp, bodyResp, errs := f.client.Client().Clone().
// 		Post(fmt.Sprintf("%s/api/v1/field/schedule", f.client.BaseURL())).
// 		Set(constants.XApiKey, apiKey).
// 		Set(constants.XserviceName, config2.Cfg.AppName).
// 		Set(constants.XrequestAt, fmt.Sprintf("%s", unixTime)).
// 		Send(string(body)).End()

// 	if len(errs) > 0 {
// 		return err
// 	}

// 	var response FieldResponse
// 	if resp.StatusCode != http.StatusCreated {
// 		err = json.Unmarshal([]byte(bodyResp), &response)
// 		if err != nil {
// 			return err
// 		}
// 		FieldError := fmt.Errorf("field response: %s", response.Message)
// 		return FieldError
// 	}

// 	err = json.Unmarshal([]byte(bodyResp), &response)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (f *FieldClient) UpdateStatus(request *dto.UpdateFieldScheduleStatusRequest) error {
	unixTime := time.Now().Unix()
	generateApikey := fmt.Sprintf("%s:%s:%d", "field-services", Cfg.Cfg.InternalService.Field.SignatureKey, unixTime)
	apiKey := util.GenerateSHA256(generateApikey)

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/field/schedule/status", Cfg.Cfg.InternalService.Field.Host), bytes.NewBuffer(body))
	//req.Header.Set("Authorization", bearerToken)
	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.XrequestAt, fmt.Sprintf("%d", unixTime))
	req.Header.Set(constants.XserviceName, "field-services")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("err: %v", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
