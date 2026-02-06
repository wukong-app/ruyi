package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wukong-app/ruyi"
	"github.com/wukong-app/ruyi/pkg/contract"
)

// main 是命令行工具的入口函数
func main() {
	kindFlag := flag.String("kind", "file", "转换类型 (file/currency/time/number)")
	fromFlag := flag.String("from", "", "源 Concept 格式 (例如 png, usd, yyyy-mm-dd)")
	toFlag := flag.String("to", "", "目标 Concept 格式 (例如 jpeg, cny, timestamp)")
	inFlag := flag.String("in", "", "输入内容: 文件路径 或 原始数据")
	outFlag := flag.String("out", "", "输出内容: 文件路径 或 原始数据输出路径")

	flag.Parse()

	if *fromFlag == "" || *toFlag == "" || *inFlag == "" || *outFlag == "" {
		fmt.Println("必须提供 from, to, in, out 参数")
		flag.Usage()
		os.Exit(1)
	}

	// 创建 Ruyi 实例
	r, err := ruyi.New()
	if err != nil {
		fmt.Printf("创建 Ruyi 实例失败: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	fromName := contract.ConceptName(*fromFlag)
	toName := contract.ConceptName(*toFlag)

	// 转换 Kind
	var kind contract.Kind
	switch strings.ToLower(*kindFlag) {
	case "file":
		kind = contract.File
	default:
		fmt.Printf("未知 kind 类型: %s\n", *kindFlag)
		os.Exit(1)
	}

	if !r.CanConvert(ctx, kind, fromName, toName) {
		fmt.Printf("不支持 %s 类型的转换: %s -> %s\n", kind, fromName, toName)
		os.Exit(1)
	}

	var outData []byte

	switch kind {
	case contract.File:
		// 文件类型: in 是文件路径
		fromData, err := os.ReadFile(*inFlag)
		if err != nil {
			fmt.Printf("读取输入文件失败: %v\n", err)
			os.Exit(1)
		}

		outData, err = r.ConvertFile(ctx, fromName, toName, fromData)
		if err != nil {
			fmt.Printf("文件转换失败: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(*outFlag, outData, 0644); err != nil {
			fmt.Printf("写入输出文件失败: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("未知 kind 类型: %s\n", *kindFlag)
		os.Exit(1)
	}

	fmt.Printf("转换成功: kind=%s, %s -> %s, 输出: %s\n", kind, fromName, toName, *outFlag)
}
