package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lin-coco/zhinao-go"
)

// downloadImage 下载图片并保存到本地
func downloadImage(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 保存图片
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

// runBasicExample 运行基础图像生成示例
func runBasicExample(ctx context.Context, client *zhinao.Client) {
	fmt.Println("1. 基础图像生成 - 蓝天白云")
	basicRequest := &zhinao.Text2ImgRequest{
		Model:  zhinao.Model360Flux1KontextDev,
		Style:  zhinao.ImageStyleRealistic,
		Prompt: "画一个蓝天白云的图片",
		Width:  512,
		Height: 512,
	}

	resp, err := client.Images.Text2Img(ctx, basicRequest)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✅ 状态: %s\n", resp.Status)
	fmt.Printf("⏱️  生成耗时: %d 秒\n", resp.GenerationTime)
	fmt.Printf("📸 生成图片数量: %d\n", len(resp.Output))

	if len(resp.Output) > 0 {
		filename := "example1_basic.png"
		if err := downloadImage(resp.Output[0], filename); err != nil {
			log.Printf("下载图片失败: %v\n", err)
		} else {
			absPath, _ := filepath.Abs(filename)
			fmt.Printf("💾 图片已保存: %s\n", absPath)
		}
	}
	fmt.Println()
}

// runNegativePromptExample 运行使用负向提示词示例
func runNegativePromptExample(ctx context.Context, client *zhinao.Client) {
	fmt.Println("2. 使用负向提示词 - 美丽风景")
	negativeRequest := &zhinao.Text2ImgRequest{
		Model:          zhinao.Model360CVC0V5,
		Style:          zhinao.ImageStyleRealistic,
		Prompt:         "美丽的山水风景，阳光明媚",
		NegativePrompt: "模糊,低质量,暗淡",
		Width:          512,
		Height:         512,
		GuidanceScale:  8.5,
	}

	resp, err := client.Images.Text2Img(ctx, negativeRequest)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✅ 状态: %s\n", resp.Status)
	fmt.Printf("📝 Prompt: %s\n", resp.Meta.Prompt)
	fmt.Printf("🚫 Negative Prompt: %s\n", resp.Meta.NegativePrompt)
	fmt.Printf("⚡ 提示词强度: %.1f\n", resp.Meta.GuidanceScale)

	if len(resp.Output) > 0 {
		filename := "example2_negative_prompt.png"
		if err := downloadImage(resp.Output[0], filename); err != nil {
			log.Printf("下载图片失败: %v\n", err)
		} else {
			absPath, _ := filepath.Abs(filename)
			fmt.Printf("💾 图片已保存: %s\n", absPath)
		}
	}
	fmt.Println()
}

// runStylesExample 运行不同风格的图像生成示例
func runStylesExample(ctx context.Context, client *zhinao.Client) {
	styles := []struct {
		style    zhinao.ImageStyle
		name     string
		filename string
	}{
		{zhinao.ImageStyleCartoon, "卡通", "example3_cartoon.png"},
		{zhinao.ImageStylePapercut, "剪纸", "example3_papercut.png"},
		{zhinao.ImageStyleCG, "CG", "example3_cg.png"},
	}

	fmt.Println("3. 不同风格的图像生成")
	for _, s := range styles {
		styleRequest := &zhinao.Text2ImgRequest{
			Model:  zhinao.ModelDoubaoSeededitV3,
			Style:  s.style,
			Prompt: "一只可爱的小猫咪",
			Width:  512,
			Height: 512,
		}

		resp, err := client.Images.Text2Img(ctx, styleRequest)
		if err != nil {
			log.Printf("  ❌ %s风格生成失败: %v\n", s.name, err)
			continue
		}

		fmt.Printf("  ✅ %s风格:\n", s.name)
		if len(resp.Output) > 0 {
			if err := downloadImage(resp.Output[0], s.filename); err != nil {
				log.Printf("     下载失败: %v\n", err)
			} else {
				absPath, _ := filepath.Abs(s.filename)
				fmt.Printf("     💾 已保存: %s\n", absPath)
			}
		}
	}
	fmt.Println()
}

// runBatchExample 运行批量生成示例
func runBatchExample(ctx context.Context, client *zhinao.Client) {
	fmt.Println("4. 批量生成多张图片")
	batchRequest := &zhinao.Text2ImgRequest{
		Model:   zhinao.ModelHunyuanImage,
		Style:   zhinao.ImageStyleRealistic,
		Prompt:  "科技感十足的未来城市",
		Width:   512,
		Height:  512,
		Samples: 3, // 一次生成3张
		Seed:    12345,
	}

	resp, err := client.Images.Text2Img(ctx, batchRequest)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✅ 状态: %s\n", resp.Status)
	fmt.Printf("⏱️  生成耗时: %d 秒\n", resp.GenerationTime)
	fmt.Printf("📸 成功生成图片数: %d\n", len(resp.Output))
	fmt.Printf("🎲 种子值: %d\n", resp.Meta.Seed)

	for i, url := range resp.Output {
		filename := fmt.Sprintf("example4_batch_%d.png", i+1)
		if err := downloadImage(url, filename); err != nil {
			log.Printf("  下载图片 %d 失败: %v\n", i+1, err)
		} else {
			absPath, _ := filepath.Abs(filename)
			fmt.Printf("  💾 图片 %d 已保存: %s\n", i+1, absPath)
		}
	}
	fmt.Println()
}

// runCustomExample 运行自定义参数示例
func runCustomExample(ctx context.Context, client *zhinao.Client) {
	fmt.Println("5. 自定义详细参数")
	customRequest := &zhinao.Text2ImgRequest{
		Model:             zhinao.ModelQwenImageEdit,
		Style:             zhinao.ImageStyleRealistic,
		Prompt:            "璀璨星空下的宁静湖泊",
		Width:             1024,
		Height:            768,
		NumInferenceSteps: 30,
		GuidanceScale:     10.0,
		Seed:              99999,
		EnhancePrompt:     true,
	}

	resp, err := client.Images.Text2Img(ctx, customRequest)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✅ 状态: %s\n", resp.Status)
	fmt.Printf("📐 图像尺寸: %dx%d\n", resp.Meta.W, resp.Meta.H)
	fmt.Printf("🔢 采样步数: %d\n", resp.Meta.Steps)
	fmt.Printf("⚡ 提示词强度: %.1f\n", resp.Meta.GuidanceScale)
	fmt.Printf("🎲 种子值: %d\n", resp.Meta.Seed)
	fmt.Printf("⏱️  生成耗时: %d 秒\n", resp.GenerationTime)

	if len(resp.Output) > 0 {
		filename := "example5_custom.png"
		if err := downloadImage(resp.Output[0], filename); err != nil {
			log.Printf("下载图片失败: %v\n", err)
		} else {
			absPath, _ := filepath.Abs(filename)
			fmt.Printf("💾 图片已保存: %s\n", absPath)
		}
	}
}

func main() {

	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	fmt.Println("=== 360智脑图像生成示例 ===")
	fmt.Println()

	// 运行各个示例
	runBasicExample(ctx, client)
	runNegativePromptExample(ctx, client)
	runStylesExample(ctx, client)
	runBatchExample(ctx, client)
	runCustomExample(ctx, client)

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("✨ 所有生成的图片已保存到当前目录")
	fmt.Println("📁 生成的文件:")
	fmt.Println("   - example1_basic.png (基础生成)")
	fmt.Println("   - example2_negative_prompt.png (负向提示词)")
	fmt.Println("   - example3_cartoon.png (卡通风格)")
	fmt.Println("   - example3_papercut.png (剪纸风格)")
	fmt.Println("   - example3_cg.png (CG风格)")
	fmt.Println("   - example4_batch_1.png, example4_batch_2.png, example4_batch_3.png (批量生成)")
	fmt.Println("   - example5_custom.png (自定义参数)")
}
