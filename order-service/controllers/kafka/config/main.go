package kafka

import (
	"order-service/config"
	kafka "order-service/controllers/kafka"
	paymentKafka "order-service/controllers/kafka/payment"

	"golang.org/x/exp/slices"
)

type Kafka struct {
	consumer *ConsumerGroup
	kafka    kafka.IKafkaRegistry
}

type IKafka interface {
	Register()
}

func NewKafkaConsumer(consumer *ConsumerGroup, kafka kafka.IKafkaRegistry) IKafka {
	return &Kafka{consumer: consumer, kafka: kafka}
}

func (k *Kafka) Register() {
	k.PaymendHandler()
}

func (k *Kafka) PaymendHandler() {
	if slices.Contains(config.Cfg.Kafka.Topic, paymentKafka.PaymentTopic) {
		k.consumer.RegisterHandler(paymentKafka.PaymentTopic, k.kafka.GetPayment().HandlePayment)
	}
}
