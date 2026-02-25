package controller

import (
	"bluebell/dao"
	"bluebell/logic"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommunityController struct {
	communityLogic logic.CommunityLogicInterface
}

func NewCommunityController(communityLogic logic.CommunityLogicInterface) *CommunityController {
	return &CommunityController{
		communityLogic: communityLogic,
	}
}

func (c *CommunityController) CommunityHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//
		list, err := c.communityLogic.GetCommunityList()
		if err != nil {
			zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
			ResponseError(ctx, CodeServerBusy) // 不轻易把服务端报错暴露给外面
			return
		}
		if len(list) == 0 {
			zap.L().Info("no community found")
			ResponseSuccess(ctx, nil)
			return
		}
		ResponseSuccess(ctx, list)
		//
	}
}

func (c *CommunityController) CommunityDetailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//
		param := ctx.Param("id")
		num, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt failed", zap.String("param", param), zap.Error(err))
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		detail, err := c.communityLogic.GetCommunityDetail(num)
		if err != nil {
			ok := errors.Is(err, dao.ErrorCommunityNotExist)
			if ok {
				zap.L().Error("community not exist", zap.Int64("id", num))
				ResponseError(ctx, CodeCommunityNotExist)
				return
			}
			zap.L().Error("logic.GetCommunityDetail() failed", zap.Int64("id", num), zap.Error(err))
			ResponseError(ctx, CodeServerBusy) // 不轻易把服务端报错暴露给外面
			return
		}
		ResponseSuccess(ctx, detail)
		//
	}
}
