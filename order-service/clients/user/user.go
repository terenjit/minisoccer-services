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
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserClient struct {
	client config.IClientConfig
}

type IUserClient interface {
	GetUserByToken(context.Context) (*UserData, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error)
}

func NewUserClient(client config.IClientConfig) IUserClient {
	return &UserClient{client: client}
}

// func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
// 	unixTime := time.Now().Unix()
// 	generateApikey := fmt.Sprintf("%s:%s:%d", config2.Cfg.AppName, u.client.SignatureKey(), unixTime)
// 	apiKey := util.GenerateSHA256(generateApikey)
// 	token := ctx.Value(constants.Token).(string)
// 	bearerToken := fmt.Sprintf("Bearer %s", token)

// 	var response UserResponse
// 	request := u.client.Client().Set(constants.Authorization, bearerToken).
// 		Set(constants.XApiKey, apiKey).
// 		Set(constants.XserviceName, config2.Cfg.AppName).
// 		Set(constants.XrequestAt, fmt.Sprintf("%s", unixTime)).
// 		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL()))

// 	resp, _, errs := request.EndStruct(&response)
// 	if len(errs) > 0 {
// 		return nil, errs[0]
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("user response: %s", response.Message)
// 	}

// 	return &response.Data, nil
// }

// func (u *UserClient) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error) {
// 	unixTime := time.Now().Unix()
// 	generateApikey := fmt.Sprintf("%s:%s:%d", config2.Cfg.AppName, u.client.SignatureKey(), unixTime)
// 	apiKey := util.GenerateSHA256(generateApikey)
// 	token := ctx.Value(constants.Token).(string)
// 	bearerToken := fmt.Sprintf("Bearer %s", token)

// 	var response UserResponse
// 	request := u.client.Client().Set(constants.Authorization, bearerToken).
// 		Set(constants.XApiKey, apiKey).
// 		Set(constants.XserviceName, config2.Cfg.AppName).
// 		Set(constants.XrequestAt, fmt.Sprintf("%s", unixTime)).
// 		Get(fmt.Sprintf("%s/api/v1/auth/%s", u.client.BaseURL(), uuid))

// 	resp, _, errs := request.EndStruct(&response)
// 	if len(errs) > 0 {
// 		return nil, errs[0]
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("user response: %s", response.Message)
// 	}

// 	return &response.Data, nil
// }

func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateApikey := fmt.Sprintf("%s:%s:%d", "user-services", u.client.SignatureKey(), unixTime)
	apiKey := util.GenerateSHA256(generateApikey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL()), nil)
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.XrequestAt, fmt.Sprintf("%d", unixTime))
	req.Header.Set(constants.XserviceName, "user-services")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("err: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var userResp UserResponse
	json.NewDecoder(resp.Body).Decode(&userResp)
	return &userResp.Data, nil
}

func (u *UserClient) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateApikey := fmt.Sprintf("%s:%s:%d", config2.Cfg.AppName, u.client.SignatureKey(), unixTime)
	apiKey := util.GenerateSHA256(generateApikey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/auth/%s", u.client.BaseURL(), uuid), nil)
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.XrequestAt, fmt.Sprintf("%d", unixTime))
	req.Header.Set(constants.XserviceName, "user-services")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("err: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var userResp UserResponse
	json.NewDecoder(resp.Body).Decode(&userResp)
	return &userResp.Data, nil
}
