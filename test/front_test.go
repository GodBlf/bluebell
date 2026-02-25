package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultBaseURL  = "http://127.0.0.1:8084"
	testTokenSecret = "godblf"
)

type testClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type frontResponse struct {
	Code int64           `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type communityInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type communityDetail struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
}

func getBaseURL() string {
	baseURL := os.Getenv("BLUEBELL_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return baseURL
}

func requireServerReady(t *testing.T, baseURL string) {
	t.Helper()
	client := newRestyClient(baseURL)
	resp, err := client.R().Get("/")
	require.NoError(t, err, "failed to access %s, please start project first", baseURL)
	require.Equal(t, http.StatusOK, resp.StatusCode(), "project is not ready at %s", baseURL)
}

func newRestyClient(baseURL string) *resty.Client {
	return resty.New().
		SetBaseURL(baseURL).
		SetTimeout(5 * time.Second)
}

func mustTestToken(t *testing.T) string {
	t.Helper()
	claims := &testClaims{
		UserID:   1,
		Username: "front_test",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			Issuer:    "front_test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testTokenSecret))
	require.NoError(t, err)
	return signed
}

func TestCommunityListFront(t *testing.T) {
	baseURL := getBaseURL()
	requireServerReady(t, baseURL)

	client := newRestyClient(baseURL)
	token := mustTestToken(t)

	var result frontResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&result).
		Get("/api/v1/community")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, int64(1000), result.Code)

	if string(result.Data) == "null" {
		return
	}

	var list []communityInfo
	require.NoError(t, json.Unmarshal(result.Data, &list))
	for _, c := range list {
		assert.NotZero(t, c.ID)
		assert.NotEmpty(t, c.Name)
	}
}

func TestCommunityDetailFront(t *testing.T) {
	baseURL := getBaseURL()
	requireServerReady(t, baseURL)

	client := newRestyClient(baseURL)
	token := mustTestToken(t)

	var listResp frontResponse
	_, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&listResp).
		Get("/api/v1/community")
	require.NoError(t, err)
	require.Equal(t, int64(1000), listResp.Code)

	targetID := int64(1)
	if string(listResp.Data) != "null" {
		var list []communityInfo
		require.NoError(t, json.Unmarshal(listResp.Data, &list))
		if len(list) > 0 {
			targetID = list[0].ID
		}
	}

	var detailResp frontResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&detailResp).
		Get(fmt.Sprintf("/api/v1/community/%d", targetID))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, int64(1000), detailResp.Code)

	var detail communityDetail
	require.NoError(t, json.Unmarshal(detailResp.Data, &detail))
	assert.Equal(t, targetID, detail.ID)
	assert.NotEmpty(t, detail.Name)
}

func TestCommunityDetailInvalidIDFront(t *testing.T) {
	baseURL := getBaseURL()
	requireServerReady(t, baseURL)

	client := newRestyClient(baseURL)
	token := mustTestToken(t)

	var result frontResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&result).
		Get("/api/v1/community/not-a-number")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, int64(1001), result.Code)
}

func TestCommunityWithoutTokenFront(t *testing.T) {
	baseURL := getBaseURL()
	requireServerReady(t, baseURL)

	client := newRestyClient(baseURL)
	var result frontResponse
	resp, err := client.R().
		SetResult(&result).
		Get("/api/v1/community")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, int64(1006), result.Code)
}
