package controller

import (
	"bluebell/dao"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UserController struct {
	userLogic logic.UserLogicInterface
}

func NewUserController(userLogic logic.UserLogicInterface) *UserController {
	return &UserController{
		userLogic: userLogic,
	}
}

func (u *UserController) SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := &models.ParamSignUp{}
		err := c.ShouldBindJSON(p)
		//check invalid
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
		//logic signup
		if err := u.userLogic.SignUp(p); err != nil {
			zap.L().Error("logic.SignUp failed", zap.Error(err))

		}

	}
}

func (u *UserController) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//check invalid
		p := &models.ParamLogin{}
		err := c.ShouldBindJSON(p)
		if err != nil {
			zap.L().Error("Login with invalid param", zap.Error(err))
			asType, ok := errors.AsType[validator.ValidationErrors](err)
			if !ok {
				ResponseError(c, CodeInvalidParam)
				return
			}
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(asType.Translate(trans)))
			return
		}
		//logic
		user, err := u.userLogic.Login(p)
		if err != nil {
			zap.L().Error("logic.Login failed", zap.Error(err))
			ok := errors.Is(err, dao.ErrorUserNotExist)
			if ok {
				ResponseError(c, CodeUserNotExist)
				return
			}
			ok = errors.Is(err, dao.ErrorInvalidPassword)
			if ok {
				ResponseError(c, CodeInvalidPassword)
				return
			}
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseSuccess(c,
			gin.H{
				"user_id":   fmt.Sprintf("%d", user.UserID),
				"token":     user.Token,
				"user_name": user.Username,
			},
		)

	}
}

func (u *UserController) SearchUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		psu := &models.ParamSearchUser{}
		err := c.ShouldBindJSON(psu)
		if err != nil {
			zap.L().Error("SearchUser with invalid param", zap.Error(err))
			asType, ok := errors.AsType[validator.ValidationErrors](err)
			if !ok {
				ResponseError(c, CodeInvalidParam)
				return
			}
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(asType.Translate(trans)))
			return
		}
		if psu.Username == "" {
			ResponseError(c, CodeInvalidParam)
			return
		}
		user, err := u.userLogic.SearchUserByName(psu.Username)
		if err != nil {
			ok := errors.Is(err, dao.ErrorUserNotExist)
			if ok {
				zap.L().Error("user exist", zap.String("username", psu.Username), zap.Error(err))
				ResponseError(c, CodeUserNotExist)
				return
			}
			zap.L().Error("logic.SearchUserByName failed", zap.String("username", psu.Username), zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
			return
			//
		}
		ResponseSuccess(c, gin.H{
			"user_id":   fmt.Sprintf("%d", user.UserID),
			"user_name": user.Username,
		})

		//
	}
}
