# TineyeLogger

Unified logging library for OpenClaw Task Router.

## Features

- 統一的記錄格式
- 支援 Go 和 Python
- Request ID 追蹤
- 非阻塞式記錄

## Installation

### Go
```bash
go get github.com/FlatMoon-Eva/TineyeLogger/pkg/logger
```

### Python
```bash
pip install git+https://github.com/FlatMoon-Eva/TineyeLogger.git
```

## Usage

### Go (ClawRouter)
```go
import "github.com/FlatMoon-Eva/TineyeLogger/pkg/logger"

// 初始化
log := logger.New("http://tineye:8889/log", logger.SourceClawRouter)

// 生成 request ID
requestID := logger.GenerateRequestID()

// 記錄收到請求
log.LogReceive(requestID, "請幫我翻譯 hello world")

// 記錄分類結果
log.LogClassify(requestID, 2, "翻譯任務", "simple-brain")

// 記錄路由
log.LogRoute(requestID, "smartllm")
```

### Python (SmartLLM)
```python
from tineyelogger import TineyeLogger, SOURCE_SMARTLLM

# 初始化
logger = TineyeLogger("http://tineye:8889/log", SOURCE_SMARTLLM)

# 從 header 取得 request ID
request_id = request.headers.get("X-Request-ID")

# 記錄執行
await logger.log_execute(
    request_id,
    model="simple-brain",
    actual_model="gemini-3-flash-preview",
    key_id="abc1234",
    input_tokens=50
)

# 記錄回應
await logger.log_response(
    request_id,
    output_tokens=20,
    latency_ms=520,
    status="success"
)
```

## Format

See [unified-logging-format.md](https://github.com/FlatMoon-Eva/flatmoon-vault/blob/main/projects/openclaw-task-router/decisions/unified-logging-format.md)
