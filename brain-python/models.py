from pydantic import BaseModel


class MarketData(BaseModel):
    symbol: str
    price: str
    timestamp: int