package middleware

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			zap.L().Info("AuthMiddleware: Authorization header is empty")
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		part := strings.SplitN(authHeader, " ", 2)
		if len(part) != 2 || part[0] != "Bearer" {
			zap.L().Info("AuthMiddleware: Authorization header format is invalid")
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		token, err := jwt.ParseJwtToken(part[1])
		if err != nil {
			ok := errors.Is(err, jwt.AuthErrorInvalid)
			if ok {
				zap.L().Info("AuthMiddleware: token is invalid")
				controller.ResponseError(c, controller.CodeInvalidToken)
				c.Abort()
				return
			}

		}

		c.Set("userID", token.UserID)
		c.Next()

		//
	}
}
