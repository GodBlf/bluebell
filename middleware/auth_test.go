package middleware

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"bluebell/settings"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type authResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func newAuthTestRouter(next gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", next)
	return r
}

func decodeAuthResponse(t *testing.T, body []byte) authResponse {
	t.Helper()
	var resp authResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func TestAuth_MissingAuthorizationHeader(t *testing.T) {
	called := false
	r := newAuthTestRouter(func(c *gin.Context) {
		called = true
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	resp := decodeAuthResponse(t, w.Body.Bytes())
	assert.Equal(t, int64(controller.CodeNeedLogin), resp.Code)
	assert.Equal(t, controller.CodeNeedLogin.Msg(), resp.Msg)
	assert.False(t, called)
}

func TestAuth_InvalidAuthorizationHeaderFormat(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{name: "invalid prefix", header: "Basic token"},
		{name: "missing token", header: "Bearer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			r := newAuthTestRouter(func(c *gin.Context) {
				called = true
				c.Status(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
			resp := decodeAuthResponse(t, w.Body.Bytes())
			assert.Equal(t, int64(controller.CodeInvalidToken), resp.Code)
			assert.Equal(t, controller.CodeInvalidToken.Msg(), resp.Msg)
			assert.False(t, called)
		})
	}
}

func TestAuth_ValidToken(t *testing.T) {
	prevCfg := settings.GlobalConfig
	t.Cleanup(func() {
		settings.GlobalConfig = prevCfg
	})
	settings.GlobalConfig = &settings.AppConfig{
		AuthConfig: &settings.AuthConfig{
			JwtExpire: 1,
		},
	}

	token, err := jwt.GenJwtToken(12345, "tester")
	require.NoError(t, err)

	r := newAuthTestRouter(func(c *gin.Context) {
		value, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "missing userID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userID": value})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		UserID int64 `json:"userID"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(12345), resp.UserID)
}

func TestAuth_MalformedTokenPanics(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer not-a-jwt-token")

	handler := Auth()
	require.Panics(t, func() {
		handler(c)
	})
}
