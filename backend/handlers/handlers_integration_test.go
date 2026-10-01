package handlers

import (
	"accode-go/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type runResponse struct {
	SessionID string `json:"sessionId"`
}

type statusResponse struct {
	Output     string `json:"output"`
	NeedsInput bool   `json:"needsInput"`
	Completed  bool   `json:"completed"`
	Error      string `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func setupRouterForTests() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/run", RunCode)
	r.GET("/status/:session_id", GetStatus)
	r.POST("/submit_input", SubmitInput)
	r.POST("/stop", StopExecution)
	return r
}

func resetSessionsForTests() {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	sessions = make(map[string]*models.Session)
	sessionID = 0
}

func performJSONRequest(t *testing.T, router *gin.Engine, method string, path string, payload any) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload failed: %v", err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeJSON[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal response failed: %v, body=%s", err, rec.Body.String())
	}
	return out
}

func TestGoRuntimeSessionFlow(t *testing.T) {
	resetSessionsForTests()
	router := setupRouterForTests()

	runRec := performJSONRequest(t, router, http.MethodPost, "/run", map[string]string{
		"code": "not-valid-accode;",
	})
	if runRec.Code != http.StatusOK {
		t.Fatalf("run endpoint failed: status=%d body=%s", runRec.Code, runRec.Body.String())
	}

	runResp := decodeJSON[runResponse](t, runRec)
	if strings.TrimSpace(runResp.SessionID) == "" {
		t.Fatalf("run endpoint returned empty session id")
	}

	var (
		combinedOutput strings.Builder
		finalStatus    statusResponse
	)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		statusRec := performJSONRequest(t, router, http.MethodGet, "/status/"+runResp.SessionID, map[string]string{})
		if statusRec.Code != http.StatusOK {
			t.Fatalf("status endpoint failed: status=%d body=%s", statusRec.Code, statusRec.Body.String())
		}

		finalStatus = decodeJSON[statusResponse](t, statusRec)
		combinedOutput.WriteString(finalStatus.Output)
		if finalStatus.Completed {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	if !finalStatus.Completed {
		t.Fatalf("session did not complete before timeout, output=%q", combinedOutput.String())
	}
	if strings.TrimSpace(combinedOutput.String()) == "" {
		t.Fatalf("expected Go runtime to emit an error for invalid code")
	}

	inputRec := performJSONRequest(t, router, http.MethodPost, "/submit_input", map[string]string{
		"sessionId": runResp.SessionID,
		"input":     "abc",
	})
	if inputRec.Code != http.StatusUnauthorized {
		t.Fatalf("submit_input should return 401 when session is not waiting input: status=%d body=%s", inputRec.Code, inputRec.Body.String())
	}

	inputResp := decodeJSON[messageResponse](t, inputRec)
	if strings.TrimSpace(inputResp.Error) == "" {
		t.Fatalf("submit_input should return an error message")
	}

	stopRec := performJSONRequest(t, router, http.MethodPost, "/stop", map[string]string{
		"sessionId": runResp.SessionID,
	})
	if stopRec.Code != http.StatusOK {
		t.Fatalf("stop endpoint failed: status=%d body=%s", stopRec.Code, stopRec.Body.String())
	}
}
