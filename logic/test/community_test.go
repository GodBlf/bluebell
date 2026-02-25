package test

import (
	"bluebell/dao"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type communityDaoStub struct {
	getListFunc   func() ([]*models.Community, error)
	getDetailFunc func(id int64) (*models.CommunityDetail, error)
}

func (c *communityDaoStub) GetCommunityList() ([]*models.Community, error) {
	if c.getListFunc != nil {
		return c.getListFunc()
	}
	return nil, nil
}

func (c *communityDaoStub) GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	if c.getDetailFunc != nil {
		return c.getDetailFunc(id)
	}
	return nil, nil
}

func TestNewCommunityLogic(t *testing.T) {
	l := logic.NewCommunityLogic(&communityDaoStub{})
	assert.NotNil(t, l)
}

func TestGetCommunityList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []*models.Community{
			{
				ID:   1,
				Name: "go",
			},
		}
		l := logic.NewCommunityLogic(&communityDaoStub{
			getListFunc: func() ([]*models.Community, error) {
				return expected, nil
			},
		})
		list, err := l.GetCommunityList()
		assert.NoError(t, err)
		assert.Equal(t, expected, list)
	})

	t.Run("dao_error", func(t *testing.T) {
		expectedErr := errors.New("query failed")
		l := logic.NewCommunityLogic(&communityDaoStub{
			getListFunc: func() ([]*models.Community, error) {
				return nil, expectedErr
			},
		})
		list, err := l.GetCommunityList()
		assert.Nil(t, list)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestGetCommunityDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const targetID int64 = 1
		expected := &models.CommunityDetail{
			ID:           targetID,
			Name:         "go",
			Introduction: "go intro",
			CreateTime:   time.Date(2026, 2, 25, 0, 0, 0, 0, time.UTC),
		}
		l := logic.NewCommunityLogic(&communityDaoStub{
			getDetailFunc: func(id int64) (*models.CommunityDetail, error) {
				assert.Equal(t, targetID, id)
				return expected, nil
			},
		})
		detail, err := l.GetCommunityDetail(targetID)
		assert.NoError(t, err)
		assert.Equal(t, expected, detail)
	})

	t.Run("community_not_exist", func(t *testing.T) {
		l := logic.NewCommunityLogic(&communityDaoStub{
			getDetailFunc: func(id int64) (*models.CommunityDetail, error) {
				return nil, dao.ErrorCommunityNotExist
			},
		})
		detail, err := l.GetCommunityDetail(100)
		assert.Nil(t, detail)
		assert.ErrorIs(t, err, dao.ErrorCommunityNotExist)
	})

	t.Run("dao_error", func(t *testing.T) {
		expectedErr := errors.New("db timeout")
		l := logic.NewCommunityLogic(&communityDaoStub{
			getDetailFunc: func(id int64) (*models.CommunityDetail, error) {
				return nil, expectedErr
			},
		})
		detail, err := l.GetCommunityDetail(100)
		assert.Nil(t, detail)
		assert.ErrorIs(t, err, expectedErr)
	})
}
