"""TineyeLogger - Unified logging for OpenClaw Task Router"""
import json
import random
import time
from datetime import datetime, timezone
from typing import Optional
import aiohttp

# Stage 定義
STAGE_RECEIVE = "receive"
STAGE_CLASSIFY = "classify"
STAGE_ROUTE = "route"
STAGE_EXECUTE = "execute"
STAGE_RESPONSE = "response"
STAGE_LOG = "log"

# Source 定義
SOURCE_CLAWROUTER = "clawrouter"
SOURCE_SMARTLLM = "smartllm"
SOURCE_TINEYE = "tineye"


def generate_request_id() -> str:
    """生成唯一 request ID"""
    now = datetime.now()
    random_hex = format(random.randint(0, 0xffffff), '06x')
    return f"req-{now:%Y%m%d-%H%M%S}-{random_hex}"


class TineyeLogger:
    """統一 logger"""
    
    def __init__(self, collector_url: str, source: str):
        self.collector_url = collector_url
        self.source = source
    
    async def log(self, record: dict) -> None:
        """記錄一筆資料"""
        # 自動填入 timestamp 和 source
        if "ts" not in record:
            record["ts"] = datetime.now(timezone.utc).isoformat()
        if "source" not in record:
            record["source"] = self.source
        
        try:
            async with aiohttp.ClientSession() as session:
                await session.post(
                    self.collector_url,
                    data=json.dumps(record).encode(),
                    headers={"Content-Type": "application/json"},
                    timeout=aiohttp.ClientTimeout(total=3)
                )
        except Exception:
            pass  # 不阻塞主流程
    
    async def log_receive(self, request_id: str, user_message: str) -> None:
        """記錄收到請求"""
        await self.log({
            "request_id": request_id,
            "stage": STAGE_RECEIVE,
            "user_message": user_message[:500]  # 限制長度
        })
    
    async def log_classify(self, request_id: str, tier: int, reason: str, model: str) -> None:
        """記錄分類結果"""
        await self.log({
            "request_id": request_id,
            "stage": STAGE_CLASSIFY,
            "tier": tier,
            "tier_reason": reason,
            "model": model
        })
    
    async def log_route(self, request_id: str, target: str) -> None:
        """記錄路由決策"""
        await self.log({
            "request_id": request_id,
            "stage": STAGE_ROUTE,
            "target": target
        })
    
    async def log_execute(
        self,
        request_id: str,
        model: str,
        actual_model: str,
        key_id: str,
        input_tokens: int
    ) -> None:
        """記錄執行開始"""
        await self.log({
            "request_id": request_id,
            "stage": STAGE_EXECUTE,
            "model": model,
            "actual_model": actual_model,
            "key_id": key_id,
            "input_tokens": input_tokens
        })
    
    async def log_response(
        self,
        request_id: str,
        output_tokens: int,
        latency_ms: int,
        status: str = "success"
    ) -> None:
        """記錄回應完成"""
        await self.log({
            "request_id": request_id,
            "stage": STAGE_RESPONSE,
            "output_tokens": output_tokens,
            "latency_ms": latency_ms,
            "status": status
        })
    
    async def log_error(self, request_id: str, error_message: str) -> None:
        """記錄錯誤"""
        await self.log({
            "request_id": request_id,
            "stage": STAGE_RESPONSE,
            "status": "error",
            "error_message": error_message
        })
