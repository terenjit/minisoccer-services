package config

import (
	"field-service/common/util"
	"os"

	"github.com/sirupsen/logrus"
)

var Cfg AppConfig

type AppConfig struct {
	Port                  int      `json:"port"`
	AppName               string   `json:"appName"`
	AppEnv                string   `json:"appEnv"`
	SignatureKey          string   `json:"signatureKey"`
	Database              Database `json:"database"`
	RateLimiterMaxRequest int      `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond int      `json:"rateLimiterTimeSecond"`
	JWTSecretKey          string   `json:"jwtSecretKey"`
	JWTExpirationTime     int      `json:"jwtExpirationTime"` // in minutes
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
