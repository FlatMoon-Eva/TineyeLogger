# TineyeLogger

Unified logging library for OpenClaw Task Router.

## Features

- Unified log format
- Support for Go and Python
- Request ID tracking
- Non-blocking logging

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

### Go Example (ClawRouter)

```go
package main

import (
    "github.com/FlatMoon-Eva/TineyeLogger/pkg/logger"
)

func main() {
    // Initialize logger
    log := logger.New(
        "http://tineye.tineye.svc.cluster.local:8889/log",
        "clawrouter",
    )

    // Generate request ID
    requestID := logger.GenerateRequestID()
    // Output: "req-20260323-150217-38ba17"

    // Log receive stage
    log.LogReceive(requestID, "請幫我翻譯 hello world")

    // Log classify stage
    tier := 2
    log.LogClassify(requestID, tier, "翻譯任務", "simple-brain")

    // Log route stage
    log.LogRoute(requestID, "smartllm")

    // Log response stage
    log.LogResponse(requestID, 20, 520, "success")

    // Custom log
    log.Log(&logger.Record{
        RequestID: requestID,
        Stage:     "custom",
        Status:    "ok",
    })
}
```

### Python Example (SmartLLM)

```python
from tineyelogger import TineyeLogger

# Initialize logger
logger = TineyeLogger(
    collector_url="http://tineye.tineye.svc.cluster.local:8889/log",
    source="smartllm"
)

# Get request ID from header
request_id = request.headers.get("X-Request-ID")

# Log execute stage
await logger.log_execute(
    request_id=request_id,
    model="simple-brain",
    actual_model="gemini-3-flash-preview",
    key_id="abc1234",
    input_tokens=50
)

# Log response stage
await logger.log_response(
    request_id=request_id,
    output_tokens=20,
    latency_ms=520,
    status="success"
)

# Custom log
await logger.log({
    "request_id": request_id,
    "stage": "custom",
    "status": "ok"
})
```

## API Reference

### Go

#### Functions

- `New(collectorURL, source string) *Logger` - Create logger instance
- `GenerateRequestID() string` - Generate unique request ID

#### Methods

- `Log(record *Record) error` - Log custom record
- `LogReceive(requestID, userMessage string) error` - Log receive stage
- `LogClassify(requestID string, tier int, reason, model string) error` - Log classify stage
- `LogRoute(requestID, target string) error` - Log route stage
- `LogExecute(requestID, model, actualModel, keyID string, inputTokens int) error` - Log execute stage
- `LogResponse(requestID string, outputTokens, latencyMs int, status string) error` - Log response stage
- `LogError(requestID, errorMessage string) error` - Log error

### Python

#### Functions

- `TineyeLogger(collector_url, source)` - Create logger instance
- `generate_request_id()` - Generate unique request ID

#### Methods

- `async log(record: dict)` - Log custom record
- `async log_receive(request_id, user_message)` - Log receive stage
- `async log_classify(request_id, tier, reason, model)` - Log classify stage
- `async log_route(request_id, target)` - Log route stage
- `async log_execute(request_id, model, actual_model, key_id, input_tokens)` - Log execute stage
- `async log_response(request_id, output_tokens, latency_ms, status)` - Log response stage
- `async log_error(request_id, error_message)` - Log error

## Format

See [unified-logging-format.md](https://github.com/FlatMoon-Eva/flatmoon-vault/blob/main/projects/openclaw-task-router/decisions/unified-logging-format.md)
