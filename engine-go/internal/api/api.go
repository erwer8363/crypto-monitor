package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendToPython(symbol string, price string) {
	postBody, _ := json.Marshal(map[string]interface{}{
		"symbol":    symbol,
		"price":     price,
		"timestamp": time.Now().Unix(),
	})
	responsBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://127.0.0.1:8000/analyze", "application/json", responsBody)
	if err != nil {
		log.Printf("❌ 发送给 Python 失败: %v", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("🚀 已将数据同步至决策大脑...")
}

func SendWhaleToPython(message string) {
	postBody, _ := json.Marshal(map[string]interface{}{
		"content": message,
	})
	http.Post("http://127.0.0.1:8000/whale_alert", "application/json", bytes.NewBuffer(postBody))
	fmt.Println("🚀 已将巨鲸数据同步至决策大脑...")
}
