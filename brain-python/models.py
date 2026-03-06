from pydantic import BaseModel


class MarketData(BaseModel):
    symbol: str
    price: float
    high: float
    low: float
    volume: float
    timestamp: int