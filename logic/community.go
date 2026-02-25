package logic

import (
	"bluebell/dao"
	"bluebell/models"
	"errors"

	"go.uber.org/zap"
)

type CommunityLogicInterface interface {
	GetCommunityList() ([]*models.Community, error)
	GetCommunityDetail(id int64) (*models.CommunityDetail, error)
}

type CommunityLogic struct {
	communityDao dao.CommunityDaoInterface
}

func NewCommunityLogic(communityDao dao.CommunityDaoInterface) *CommunityLogic {
	return &CommunityLogic{communityDao: communityDao}
}

func (c *CommunityLogic) GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	detail, err := c.communityDao.GetCommunityDetail(id)
	if err != nil {
		ok := errors.Is(err, dao.ErrorCommunityNotExist)
		if ok {
			zap.L().Error("community not exist", zap.Int64("id", id))
			return nil, dao.ErrorCommunityNotExist
		}
		zap.L().Error("query community detail failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return detail, nil
}

func (c *CommunityLogic) GetCommunityList() ([]*models.Community, error) {
	return c.communityDao.GetCommunityList()
}
