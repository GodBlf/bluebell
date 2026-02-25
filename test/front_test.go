package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultBaseURL   = "http://127.0.0.1:8084"
	testTokenSecret  = "godblf"
	defaultStartWait = 20 * time.Second
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

var startedCmd *exec.Cmd

func TestMain(m *testing.M) {
	baseURL := getBaseURL()

	if !serverReady(baseURL) {
		cmd, err := startProjectServer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to start project: %v\n", err)
			os.Exit(1)
		}
		startedCmd = cmd

		if err := waitServerReady(baseURL, defaultStartWait); err != nil {
			_ = stopProjectServer(startedCmd)
			fmt.Fprintf(os.Stderr, "server not ready: %v\n", err)
			os.Exit(1)
		}
	}

	exitCode := m.Run()
	_ = stopProjectServer(startedCmd)
	os.Exit(exitCode)
}

func getBaseURL() string {
	baseURL := os.Getenv("BLUEBELL_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return baseURL
}

func projectRootDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	if fileExists(filepath.Join(wd, "go.mod")) {
		return wd
	}
	parent := filepath.Dir(wd)
	if fileExists(filepath.Join(parent, "go.mod")) {
		return parent
	}
	return wd
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func startProjectServer() (*exec.Cmd, error) {
	cmd := exec.Command("go", "run", "./cmd/main.go")
	cmd.Dir = projectRootDir()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func stopProjectServer(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	_ = cmd.Process.Kill()
	_, err := cmd.Process.Wait()
	return err
}

func serverReady(baseURL string) bool {
	client := newRestyClient(baseURL)
	resp, err := client.R().Get("/")
	if err != nil {
		return false
	}
	return resp.StatusCode() == http.StatusOK
}

func waitServerReady(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if serverReady(baseURL) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return errors.New("timeout waiting for server")
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
	client := newRestyClient(getBaseURL())
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
	client := newRestyClient(getBaseURL())
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
	client := newRestyClient(getBaseURL())
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
	client := newRestyClient(getBaseURL())
	var result frontResponse
	resp, err := client.R().
		SetResult(&result).
		Get("/api/v1/community")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, int64(1006), result.Code)
}
