import os
import json
from google import genai
from dotenv import load_dotenv

load_dotenv()

def ask_ai_decision(symbol, price, rsi, whale_info):
    client = genai.Client(api_key=os.getenv("GEMINI_API_KEY"))

    prompt = f"""
    你是我的私人量化专家。目前账户本金 1000 USDT，目标是稳健复利。
    
    【当前市场行情】
    - 交易对: {symbol}
    - 当前价格: {price}
    - 技术指标：RSI(14)为 {rsi:.2f}
    - 链上巨鲸雷达: {whale_info}
    
    【你的任务】
    请结合 RSI 的超买超卖情况和大户的资金流向，给出决策。
    注意：如果 RSI 超卖且大户在提币，是强烈的买入信号；如果 RSI 超买且大户在充币，是强烈卖出信号。
    
    请严格返回如下 JSON 格式：
    {{"action": "BUY/SELL/HOLD", "confidence": 1-100, "reason": "一句话理由"}}
    """

    try:
        response = client.models.generate_content(
            model="gemini-1.5-flash",
            contents= prompt
        )
        clear_json = response.text.strip().replace("```json", "").replace("```", "")
        return json.loads(clear_json)
    except Exception as e:
        print(f"❌ AI 调用失败: {e}")
        return '{"action": "HOLD", "reason": "AI Error"}'