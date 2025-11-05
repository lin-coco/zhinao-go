package main

import (
	"context"
	"fmt"
	"github.com/lin-coco/zhinao-go"
	"log"
)

func main() {
	// 创建客户端（从环境变量读取 API Key）
	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 获取可用模型列表 ===")

	// 获取模型列表
	models, err := client.ListModels(context.Background())
	if err != nil {
		log.Fatalf("获取模型列表失败: %v\n", err)
	}

	fmt.Printf("找到 %d 个可用模型:\n\n", len(models.Data))

	// 打印模型列表
	for i, model := range models.Data {
		fmt.Printf("%d. %s\n", i+1, model.ID)
		fmt.Printf("   类型: %s\n", model.Object)
		fmt.Printf("   所有者: %s\n", model.OwnedBy)
		fmt.Println()
	}

	// 获取特定模型的详细信息
	if len(models.Data) > 0 {
		modelID := models.Data[0].ID
		fmt.Printf("=== 获取模型 '%s' 的详细信息 ===\n\n", modelID)

		info, err := client.GetModel(context.Background(), modelID)
		if err != nil {
			log.Fatalf("获取模型信息失败: %v\n", err)
		}

		fmt.Printf("模型 ID: %s\n", info.ID)
		fmt.Printf("类型: %s\n", info.Object)
		fmt.Printf("所有者: %s\n", info.OwnedBy)
	}

	// 也可以直接查询指定模型
	fmt.Println("\n=== 查询指定模型 ===")

	turboInfo, err := client.GetModel(context.Background(), zhinao.Model360GPTTurbo)
	if err != nil {
		log.Printf("获取 %s 模型失败: %v\n", zhinao.Model360GPTTurbo, err)
	} else {
		fmt.Printf("✓ %s 可用\n", turboInfo.ID)
		fmt.Printf("  所有者: %s\n", turboInfo.OwnedBy)
	}
}
