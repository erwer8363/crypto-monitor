package main

import (
	"context"
	"engine-go/internal/api"
	"engine-go/services"
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance/v2"
)

func main() {
	wssAddr := "wss://eth-mainnet.g.alchemy.com/v2/W3hiLRpFiTSmjGwT-et_y"

	go services.WatchWhales(wssAddr, api.SendWhaleToPython)

	client := binance.NewClient("", "")
	fmt.Println("🚀 现货监控系统启动...")
	for {
		symbol := "BTCUSDT"
		prices, err := client.NewListPricesService().Symbol(symbol).Do(context.Background())
		if err != nil {
			log.Printf("❌ 获取价格失败: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		for _, price := range prices {
			fmt.Printf("⏰ 时间: %s | 币种: %s | 当前价格: %s\n", time.Now().Format("15:04:05"), symbol, price.Price)
			api.SendToPython(symbol, price.Price)
		}
		time.Sleep(60 * time.Second)
	}
	select {}
}
