package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lin-coco/zhinao-go"
)

func main() {
	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	fmt.Println("=== 360智脑向量生成示例 ===")
	fmt.Println()

	// 示例 1: 基础向量生成
	fmt.Println("1. 基础向量生成")
	basicReq := &zhinao.EmbeddingsRequest{
		Model: zhinao.ModelEmbeddingS1V1,
		Input: []string{"你好"},
	}

	resp, err := client.Embeddings.Create(ctx, basicReq)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✅ 模型: %s\n", resp.Model)
		fmt.Printf("📊 生成向量数: %d\n", len(resp.Data))
		fmt.Printf("📏 向量维度: %d\n", len(resp.Data[0].Embedding))
		fmt.Printf("🔢 Token 消耗: %d (总计: %d)\n", resp.Usage.PromptTokens, resp.Usage.TotalTokens)
		fmt.Printf("🎯 向量前5维: %.6f, %.6f, %.6f, %.6f, %.6f\n",
			resp.Data[0].Embedding[0],
			resp.Data[0].Embedding[1],
			resp.Data[0].Embedding[2],
			resp.Data[0].Embedding[3],
			resp.Data[0].Embedding[4])
		fmt.Println()
	}

	// 示例 2: 批量生成向量
	fmt.Println("2. 批量生成向量")
	batchReq := &zhinao.EmbeddingsRequest{
		Model: zhinao.ModelEmbeddingS1V1,
		Input: []string{
			"你好，世界",
			"人工智能",
			"自然语言处理",
		},
	}

	resp, err = client.Embeddings.Create(ctx, batchReq)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✅ 模型: %s\n", resp.Model)
		fmt.Printf("📊 生成向量数: %d\n", len(resp.Data))
		fmt.Printf("🔢 Token 消耗: %d (总计: %d)\n", resp.Usage.PromptTokens, resp.Usage.TotalTokens)
		for i, data := range resp.Data {
			fmt.Printf("  向量 %d: 维度=%d, 索引=%d\n", i+1, len(data.Embedding), data.Index)
		}
		fmt.Println()
	}

	// 示例 3: 使用 Builder 构建请求
	fmt.Println("3. 使用 Builder 构建请求")
	builderReq := zhinao.NewEmbeddings(zhinao.ModelEmbeddingS1V1).
		AddInput("机器学习").
		AddInput("深度学习").
		AddInput("神经网络").
		SetUser("example-user").
		Build()

	resp, err = client.Embeddings.Create(ctx, builderReq)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✅ 模型: %s\n", resp.Model)
		fmt.Printf("📊 生成向量数: %d\n", len(resp.Data))
		fmt.Printf("🔢 Token 消耗: %d\n", resp.Usage.TotalTokens)
		fmt.Println()
	}

	fmt.Println("=== 示例完成 ===")
	fmt.Println("💡 提示:")
	fmt.Println("   - 向量可用于语义搜索、文本分类、聚类等任务")
	fmt.Println("   - 批量生成向量可以提高效率")
	fmt.Println("   - 向量相似度计算建议使用专业向量数据库（如 Milvus、Pinecone）")
	fmt.Println("   - 或使用第三方数学库（如 gonum.org/v1/gonum）")
}
