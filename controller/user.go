package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UserController struct {
	userLogic *logic.UserLogic
}

func NewUserController(userLogic *logic.UserLogic) *UserController {
	return &UserController{
		userLogic: userLogic,
	}
}

func (u *UserController) SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := &models.ParamSignUp{}
		err := c.ShouldBindJSON(p)
		if err != nil {
			zap.L().Error("SignUp with invalid param", zap.Error(err))
			asType, ok := errors.AsType[validator.ValidationErrors](err)
			if !ok {
				ResponseError(c, CodeInvalidParam)
				return
			}
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(asType.Translate(trans)))
			return
		}
		if err := u.userLogic.SignUp(p); err != nil {
			zap.L().Error("logic.SignUp failed", zap.Error(err))

		}

	}
}
