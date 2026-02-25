package dao

import (
	"bluebell/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type CommunityDaoInterface interface {
	GetCommunityList() ([]*models.Community, error)
	GetCommunityDetail(id int64) (*models.CommunityDetail, error)
}

type CommunityDao struct {
	client *sqlx.DB
}

func NewCommunityDao(client *sqlx.DB) *CommunityDao {
	return &CommunityDao{client: client}
}

func (c *CommunityDao) GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	com := &models.CommunityDetail{}
	err := c.client.Get(com, "select community_id, community_name, introduction, create_time from community where community_id = ? ", id)
	if err != nil {
		ok := errors.Is(err, sql.ErrNoRows)
		if ok {
			zap.L().Warn("there is no community in db")
			return nil, ErrorCommunityNotExist
		}
		zap.L().Error("query community detail failed", zap.Error(err))
		return nil, err
	}
	return com, nil
}

func (c *CommunityDao) GetCommunityList() ([]*models.Community, error) {
	communities := make([]*models.Community, 0, 2)
	err := c.client.Select(&communities, "select community_id, community_name from community")
	if err != nil {
		ok := errors.Is(err, sql.ErrNoRows)
		if ok {
			zap.L().Warn("there is no community in db")
			return nil, nil
		}
		zap.L().Error("query community list failed", zap.Error(err))
		return nil, err
	}
	return communities, nil
}
