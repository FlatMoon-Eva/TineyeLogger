package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Stage 定義
const (
	StageReceive  = "receive"
	StageClassify = "classify"
	StageRoute    = "route"
	StageExecute  = "execute"
	StageResponse = "response"
	StageLog      = "log"
)

// Source 定義
const (
	SourceClawRouter = "clawrouter"
	SourceSmartLLM   = "smartllm"
	SourceTineye     = "tineye"
)

// Record 統一記錄格式
type Record struct {
	RequestID    string `json:"request_id"`
	Timestamp    string `json:"ts"`
	Source       string `json:"source"`
	Stage        string `json:"stage"`
	
	// 選填欄位
	UserMessage  string `json:"user_message,omitempty"`
	Tier         *int   `json:"tier,omitempty"`
	TierReason   string `json:"tier_reason,omitempty"`
	Model        string `json:"model,omitempty"`
	ActualModel  string `json:"actual_model,omitempty"`
	KeyID        string `json:"key_id,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
	LatencyMs    int    `json:"latency_ms,omitempty"`
	Status       string `json:"status,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	FallbackTo   string `json:"fallback_to,omitempty"`
	Target       string `json:"target,omitempty"`
}

// Logger 實例
type Logger struct {
	collectorURL string
	source       string
	client       *http.Client
}

// New 建立新的 logger
func New(collectorURL, source string) *Logger {
	return &Logger{
		collectorURL: collectorURL,
		source:       source,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// GenerateRequestID 生成唯一 request ID
func GenerateRequestID() string {
	now := time.Now()
	random := fmt.Sprintf("%06x", rand.Intn(0xffffff))
	return fmt.Sprintf("req-%s-%s",
		now.Format("20060102-150405"),
		random)
}

// Log 記錄一筆資料
func (l *Logger) Log(record *Record) error {
	// 自動填入 timestamp 和 source
	if record.Timestamp == "" {
		record.Timestamp = time.Now().Format(time.RFC3339Nano)
	}
	if record.Source == "" {
		record.Source = l.source
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", l.collectorURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("collector returned %d", resp.StatusCode)
	}

	return nil
}

// LogReceive 記錄收到請求
func (l *Logger) LogReceive(requestID, userMessage string) error {
	return l.Log(&Record{
		RequestID:   requestID,
		Stage:       StageReceive,
		UserMessage: userMessage,
	})
}

// LogClassify 記錄分類結果
func (l *Logger) LogClassify(requestID string, tier int, reason, model string) error {
	return l.Log(&Record{
		RequestID:  requestID,
		Stage:      StageClassify,
		Tier:       &tier,
		TierReason: reason,
		Model:      model,
	})
}

// LogRoute 記錄路由決策
func (l *Logger) LogRoute(requestID, target string) error {
	return l.Log(&Record{
		RequestID: requestID,
		Stage:     StageRoute,
		Target:    target,
	})
}

// LogExecute 記錄執行開始
func (l *Logger) LogExecute(requestID, model, actualModel, keyID string, inputTokens int) error {
	return l.Log(&Record{
		RequestID:   requestID,
		Stage:       StageExecute,
		Model:       model,
		ActualModel: actualModel,
		KeyID:       keyID,
		InputTokens: inputTokens,
	})
}

// LogResponse 記錄回應完成
func (l *Logger) LogResponse(requestID string, outputTokens, latencyMs int, status string) error {
	return l.Log(&Record{
		RequestID:    requestID,
		Stage:        StageResponse,
		OutputTokens: outputTokens,
		LatencyMs:    latencyMs,
		Status:       status,
	})
}

// LogError 記錄錯誤
func (l *Logger) LogError(requestID, errorMessage string) error {
	return l.Log(&Record{
		RequestID:    requestID,
		Stage:        StageResponse,
		Status:       "error",
		ErrorMessage: errorMessage,
	})
}
