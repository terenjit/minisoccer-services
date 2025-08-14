package config

import (
	"os"
	"payment-service/common/util"

	"github.com/sirupsen/logrus"
)

var Cfg AppConfig

type AppConfig struct {
	Port                       int             `json:"port"`
	AppName                    string          `json:"appName"`
	AppEnv                     string          `json:"appEnv"`
	SignatureKey               string          `json:"signatureKey"`
	Database                   Database        `json:"database"`
	RateLimiterMaxRequest      int             `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond      int             `json:"rateLimiterTimeSecond"`
	InternalService            InternalService `json:"internalService"`
	GcsType                    string          `json:"gcsType"`
	GcsProjectID               string          `json:"gcsProjectID"`
	GcsPrivateKeyID            string          `json:"gcsPrivateKeyID"`
	GcsPrivateKey              string          `json:"gcsPrivateKey"`
	GcsClientEmail             string          `json:"gcsClientEmail"`
	GcsClientID                string          `json:"gcsClientID"`
	GcsAuthURI                 string          `json:"gcsAuthURI"`
	GcsTokenURI                string          `json:"gcsTokenURI"`
	GcsAuthProviderX509CertURL string          `json:"gcsAuthProviderX509CertURL"`
	GcsClientX509CertURL       string          `json:"gcsClientX509CertURL"`
	GcsUniverseDomain          string          `json:"gcsUniverseDomain"`
	GcsBucketName              string          `json:"gcsBucketName"`
	Kafka                      Kafka           `json:"kafka"`
	Midtrans                   Midtrans        `json:"midtrans"`
}

type Database struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnection     int    `json:"maxOpenConnection"`
	MaxLifetimeConnection int    `json:"maxLifetimeConnection"`
	MaxIdleConnection     int    `json:"maxIdleConnection"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

type InternalService struct {
	User User `json:"user"`
}

type User struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

type Kafka struct {
	Brokers     []string `json:"brokers"`
	TimeoutInMs int      `json:"timeoutInMs"`
	MaxRetry    int      `json:"maxRetry"`
	Topic       string   `json:"topic"`
}

type Midtrans struct {
	ServerKey    string `json:"clientKey"`
	ClientKey    string `json:"secretKey"`
	IsProduction bool   `json:"isProduction"`
}

func Init() {
	err := util.BindFromJSON(&Cfg, "config.json", ".")
	if err != nil {
		logrus.Infof("failed to bind config: %v", err)
		err = util.BindFromConsul(&Cfg, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
