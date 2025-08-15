package main

import (
	"fmt"
	"log"

	"simplied-bitcoin-network-go/pkg/utils"
)

func main() {
	config, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	fmt.Printf("应用名称: %s\n", config.App.Name)
	fmt.Printf("应用版本: %s\n", config.App.Version)
	fmt.Printf("网络端口: %d\n", config.Network.Port)
	fmt.Printf("RPC端口: %d\n", config.RPC.Port)
	fmt.Printf("数据目录: %s\n", config.Blockchain.DataDir)
}
