package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"payment-service/clients"
	midtransCLient "payment-service/clients/midtrans"
	"payment-service/common/gcs"
	"payment-service/common/response"
	"payment-service/config"
	"payment-service/constants"
	controllers "payment-service/controllers/http"
	"payment-service/controllers/kafka"
	"payment-service/domain/models"
	"payment-service/middlewares"
	"payment-service/repositories"
	"payment-service/routes"
	"payment-service/services"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "start the server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
		config.Init()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		err = db.AutoMigrate(
			&models.Payment{},
			&models.PaymentHistory{},
		)
		if err != nil {
			panic(err)
		}
		gcs := initGCS()
		kafka := kafka.NewKafkaRegistry(config.Cfg.Kafka.Brokers)
		midtrans := midtransCLient.NewMidtransClient(config.Cfg.Midtrans.ServerKey, config.Cfg.Midtrans.IsProduction)
		client := clients.NewClientRegistry()
		repositories := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repositories, gcs, kafka, midtrans)
		controller := controllers.NewControllerRegistry(service)

		router := gin.Default()
		router.Use(middlewares.HandlePanic())
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to Payment services",
			})
		})

		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-request-at, x-api-key")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		lmt := tollbooth.NewLimiter(
			float64(config.Cfg.RateLimiterMaxRequest),
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Cfg.RateLimiterTimeSecond) * time.Second,
			})
		router.Use(middlewares.RateLimiter(lmt))

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group, client)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Cfg.Port)
		router.Run(port)
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}

func initGCS() gcs.IGCSClient {
	decode, err := base64.StdEncoding.DecodeString(config.Cfg.GcsPrivateKey)
	if err != nil {
		panic(err)
	}

	stringPrivateKey := string(decode)
	gcsServiceAccount := gcs.ServiceAccountKeyJSON{
		Type:                    config.Cfg.GcsType,
		ProjectID:               config.Cfg.GcsProjectID,
		PrivateKeyID:            stringPrivateKey,
		ClientEmail:             config.Cfg.GcsClientEmail,
		ClientID:                config.Cfg.GcsClientID,
		AuthURI:                 config.Cfg.GcsAuthURI,
		TokenURI:                config.Cfg.GcsTokenURI,
		AuthProviderX509CertURL: config.Cfg.GcsAuthProviderX509CertURL,
		ClientX509CertURL:       config.Cfg.GcsClientX509CertURL,
		UniverseDomain:          config.Cfg.GcsUniverseDomain,
	}
	gcsClient := gcs.NewGCSClient(
		gcsServiceAccount,
		config.Cfg.GcsBucketName,
	)
	return gcsClient
}
