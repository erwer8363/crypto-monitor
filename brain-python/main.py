import pandas as pd
import uvicorn
import pandas_ta as ta
from fastapi import FastAPI

from ai_agent import ask_ai_decision
from models import MarketData

app = FastAPI()

price_history = []
# 定义全局状态，缓存最近的情报
latest_whale_info = "市场平稳，无巨鲸异动"

@app.post("/analyze")
async def receive_data(data:MarketData):
    global price_history, latest_whale_info

    price_history.append(data.price)
    if len(price_history) > 30:
        price_history.pop(0)

    if len(price_history) < 14:
        return {"action": "COLLECTING_DATA", "count": len(price_history)}

    df = pd.DataFrame(price_history, columns=["close"])
    df["close"] = df["close"].astype(float)
    rsi = ta.rsi(df["close"], length=14).iloc[-1]
    should_call_ai = False

    if rsi < 35 or rsi > 65:
        print("🤖 触发 AI 研判...")
        should_call_ai = True

    elif "捕获巨鲸" in latest_whale_info:
        should_call_ai = True

    if should_call_ai:
        print("🤖 达到触发阈值，向 Gemini 请教...")
        decision = ask_ai_decision(data.symbol, data.price, rsi, latest_whale_info)
        latest_whale_info = "市场平稳，无巨鲸异动"
        return decision

    return {
        "action": "HOLD",
        "rsi": round(rsi, 2),
        "info": "指标正常，继续观察"
    }

@app.post("/whale_alert")
async def receive_whale_data(msg: dict):
    global latest_whale_info
    # 这里可以存入 Redis 或全局变量，供下一次 AI 分析时引用
    latest_whale_info = msg.get("content", "解析失败的巨鲸情报")
    print(f"🐋 捕获巨鲸情报: {latest_whale_info}")
    return {"status": "saved"}

# 按装订区域中的绿色按钮以运行脚本。
if __name__ == '__main__':
    uvicorn.run(app, host="127.0.0.1", port=8000)
