package services

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WatchWhales(wssAddr string, callback func(string)) {
	client, err := ethclient.Dial(wssAddr)
	if err != nil {
		log.Fatal("❌ 无法连接到节点:", err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal("❌ 订阅失败:", err)
	}
	fmt.Println("🛰️ 链上雷达已启动，正在扫描巨鲸...")
	for {
		select {
		case err := <-sub.Err():
			log.Println("⚠️ 订阅错误:", err)
		case header := <-headers:
			// 获取区块内容
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				continue
			}

			// 设定阈值：比如 100 ETH
			threshold := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))

			// 遍历区块中的所有交易
			for _, tx := range block.Transactions() {
				if tx.Value().Cmp(threshold) >= 0 {
					// 转换为 ETH 单位
					ethVal := new(big.Float).Quo(new(big.Float).SetInt(tx.Value()), big.NewFloat(1e18))
					msg := fmt.Sprintf("🐋 捕获巨鲸！金额: %.2f ETH | Hash: %s\n", ethVal, tx.Hash().Hex())
					callback(msg)
				}
			}
		}
	}
}
