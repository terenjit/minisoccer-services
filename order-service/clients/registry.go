package clients

import (
	"order-service/clients/config"
	clientsField "order-service/clients/field"
	clientsPayment "order-service/clients/payment"
	clientsUser "order-service/clients/user"
	config2 "order-service/config"
)

type ClientRegistry struct{}

type IClientRegistry interface {
	GetUser() clientsUser.IUserClient
	GetPayment() clientsPayment.IPaymentClient
	GetField() clientsField.IFieldClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}

func (c *ClientRegistry) GetUser() clientsUser.IUserClient {
	return clientsUser.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(config2.Cfg.InternalService.User.Host), config.WithSignatureKey(config2.Cfg.InternalService.User.SignatureKey)))
}

func (c *ClientRegistry) GetPayment() clientsPayment.IPaymentClient {
	return clientsPayment.NewPaymentClient(
		config.NewClientConfig(
			config.WithBaseURL(config2.Cfg.InternalService.User.Host), config.WithSignatureKey(config2.Cfg.InternalService.User.SignatureKey)))
}

func (c *ClientRegistry) GetField() clientsField.IFieldClient {
	return clientsField.NewFieldClient(
		config.NewClientConfig(
			config.WithBaseURL(config2.Cfg.InternalService.User.Host), config.WithSignatureKey(config2.Cfg.InternalService.User.SignatureKey)))
}
