package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"field-service/clients"
	"field-service/common/response"
	"field-service/config"
	"field-service/constants"
	errConstant "field-service/constants/error"
	"fmt"
	"net/http"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("recovered from panic: %v", r)
				ctx.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError.Error(),
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if httpError != nil {
			ctx.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstant.ErrToManyRequests.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

func validateAPIKEY(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XrequestAt)
	serviceName := c.GetHeader(constants.XserviceName)

	var signatureKey string
	switch serviceName {
	case "field-services":
		signatureKey = config.Cfg.SignatureKey
	case "user-services":
		signatureKey = config.Cfg.InternalService.User.SignatureKey
	default:
		return errConstant.ErrUnauthorized
	}

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		return errConstant.ErrUnauthorized
	}
	return nil
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func CheckRole(roles []string, client clients.IClientRegistry) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := client.GetUser().GetUserByToken(ctx.Request.Context())
		if err != nil {
			responseUnauthorized(ctx, errConstant.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.Role) {
			responseUnauthorized(ctx, errConstant.ErrUnauthorized.Error())
			return
		}
		ctx.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		token := ctx.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(ctx, errConstant.ErrUnauthorized.Error())
			return
		}

		err = validateAPIKEY(ctx)
		if err != nil {
			responseUnauthorized(ctx, err.Error())
			return
		}

		tokenString := extractBearerToken(token)
		tokenUser := ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), constants.Token, tokenString))
		ctx.Request = tokenUser
		ctx.Next()
	}
}

func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		err := validateAPIKEY(ctx)
		if err != nil {
			responseUnauthorized(ctx, err.Error())
			return
		}

		ctx.Next()
	}
}
