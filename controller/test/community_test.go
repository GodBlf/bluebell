package test

import (
	"bluebell/controller"
	"bluebell/dao"
	"bluebell/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type communityLogicStub struct {
	getListFunc   func() ([]*models.Community, error)
	getDetailFunc func(id int64) (*models.CommunityDetail, error)
}

func (c *communityLogicStub) GetCommunityList() ([]*models.Community, error) {
	if c.getListFunc != nil {
		return c.getListFunc()
	}
	return nil, nil
}

func (c *communityLogicStub) GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	if c.getDetailFunc != nil {
		return c.getDetailFunc(id)
	}
	return nil, nil
}

func decodeResponseData(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	resp := make(map[string]any)
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	return resp
}

func TestNewCommunityController(t *testing.T) {
	ctrl := controller.NewCommunityController(&communityLogicStub{})
	assert.NotNil(t, ctrl)
}

func TestCommunityHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("logic_error", func(t *testing.T) {
		ctrl := controller.NewCommunityController(&communityLogicStub{
			getListFunc: func() ([]*models.Community, error) {
				return nil, errors.New("db down")
			},
		})
		r := gin.New()
		r.GET("/community", ctrl.CommunityHandler())

		req := httptest.NewRequest(http.MethodGet, "/community", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponseData(t, rr)
		assert.Equal(t, float64(controller.CodeServerBusy), resp["code"])
		assert.Equal(t, controller.CodeServerBusy.Msg(), resp["msg"])
		assert.Nil(t, resp["data"])
	})

	t.Run("empty_list", func(t *testing.T) {
		ctrl := controller.NewCommunityController(&communityLogicStub{
			getListFunc: func() ([]*models.Community, error) {
				return []*models.Community{}, nil
			},
		})
		r := gin.New()
		r.GET("/community", ctrl.CommunityHandler())

		req := httptest.NewRequest(http.MethodGet, "/community", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponseData(t, rr)
		assert.Equal(t, float64(controller.CodeSuccess), resp["code"])
		assert.Equal(t, controller.CodeSuccess.Msg(), resp["msg"])
		assert.Nil(t, resp["data"])
	})

	t.Run("non_empty_list", func(t *testing.T) {
		ctrl := controller.NewCommunityController(&communityLogicStub{
			getListFunc: func() ([]*models.Community, error) {
				return []*models.Community{
					{
						ID:   1,
						Name: "golang",
					},
				}, nil
			},
		})
		r := gin.New()
		r.GET("/community", ctrl.CommunityHandler())

		req := httptest.NewRequest(http.MethodGet, "/community", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, rr.Body.String(), "non-empty list should return a JSON response")
		resp := decodeResponseData(t, rr)
		assert.Equal(t, float64(controller.CodeSuccess), resp["code"])
	})
}

func TestCommunityDetailHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid_id", func(t *testing.T) {
		ctrl := controller.NewCommunityController(&communityLogicStub{})
		r := gin.New()
		r.GET("/community/:id", ctrl.CommunityDetailHandler())

		req := httptest.NewRequest(http.MethodGet, "/community/not-num", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponseData(t, rr)
		assert.Equal(t, float64(controller.CodeInvalidParam), resp["code"])
	})

	t.Run("community_not_exist", func(t *testing.T) {
		ctrl := controller.NewCommunityController(&communityLogicStub{
			getDetailFunc: func(id int64) (*models.CommunityDetail, error) {
				return nil, dao.ErrorCommunityNotExist
			},
		})
		r := gin.New()
		r.GET("/community/:id", ctrl.CommunityDetailHandler())

		req := httptest.NewRequest(http.MethodGet, "/community/2", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponseData(t, rr)
		assert.Equal(t, float64(controller.CodeCommunityNotExist), resp["code"])
	})

	t.Run("logic_error", func(t *testing.T) {
		ctrl := controller.NewCommunityController(&communityLogicStub{
			getDetailFunc: func(id int64) (*models.CommunityDetail, error) {
				return nil, errors.New("db timeout")
			},
		})
		r := gin.New()
		r.GET("/community/:id", ctrl.CommunityDetailHandler())

		req := httptest.NewRequest(http.MethodGet, "/community/2", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponseData(t, rr)
		assert.Equal(t, float64(controller.CodeServerBusy), resp["code"])
	})

	t.Run("success", func(t *testing.T) {
		const expectedID int64 = 3
		ctrl := controller.NewCommunityController(&communityLogicStub{
			getDetailFunc: func(id int64) (*models.CommunityDetail, error) {
				assert.Equal(t, expectedID, id)
				return &models.CommunityDetail{
					ID:           id,
					Name:         "go",
					Introduction: "go community",
					CreateTime:   time.Date(2026, 2, 25, 0, 0, 0, 0, time.UTC),
				}, nil
			},
		})
		r := gin.New()
		r.GET("/community/:id", ctrl.CommunityDetailHandler())

		req := httptest.NewRequest(http.MethodGet, "/community/3", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponseData(t, rr)
		assert.Equal(t, float64(controller.CodeSuccess), resp["code"])
		data, ok := resp["data"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, float64(expectedID), data["id"])
		assert.Equal(t, "go", data["name"])
	})
}
