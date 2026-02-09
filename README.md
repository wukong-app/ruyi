# 如意 (Ruyi)

> 如意金箍棒，随心所欲。Ruyi 能将万物转化为你所需的形态。

Ruyi 是一个基于 Go 语言开发的通用格式转换工具库。它采用插件化架构，旨在为 `wukong` 项目提供强大、灵活且可扩展的数据转换核心。

## 🌟 核心特性

- **插件化架构**: 转换逻辑封装为独立的 Converter，易于扩展和维护。
- **类型安全**: 基于 Concept（概念）和 Kind（类型）的强类型设计。
- **参数化控制**: 支持在转换过程中传递参数（如图片缩放、质量控制）。
- **统一接口**: 通过统一的 Registry 和 Engine 进行管理，屏蔽底层差异。

## 📦 安装

```bash
go get github.com/wukong-app/ruyi
```

## 🚀 快速开始

### 1. 初始化引擎

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wukong-app/ruyi"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func main() {
	// 初始化 Ruyi 引擎
	ry, err := ruyi.New()
	if err != nil {
		log.Fatalf("failed to create ruyi: %v", err)
	}

	// 示例数据（实际使用时应为真实文件字节）
	inputBytes := []byte("fake-jpeg-data") 
	ctx := context.Background()

	// ... 接下文
}
```

### 2. 执行转换

```go
	// 获取 JPEG 到 PNG 的转换器
	// contract.File 表示转换类型为文件
	// contract.JPEG 和 contract.PNG 分别表示源和目标格式
	converter, err := ry.GetConverter(ctx, contract.File, contract.JPEG, contract.PNG)
	if err != nil {
		log.Fatalf("converter not found: %v", err)
	}

	// 准备参数（可选）
	params := map[string]string{
		"width":  "200", // 目标宽度
		"height": "200", // 目标高度
	}

	// 执行转换
	outputBytes, err := converter.Convert(ctx, inputBytes, params)
	if err != nil {
		// 注意：这里的 inputBytes 是伪造的，实际运行会报错，需使用真实图片数据
		log.Printf("conversion failed (expected for fake data): %v", err)
	} else {
		fmt.Printf("Successfully converted %d bytes to %d bytes\n", len(inputBytes), len(outputBytes))
	}
```

## 🔌 支持的转换

目前 Ruyi 主要支持以下图片格式的转换：

| 源格式  | 目标格式 | 支持参数说明                                                                               |
| :----- | :------ |:-------------------------------------------------------------------------------------|
| **JPEG** | **PNG**  | - `width`, `height`：支持调整输出图片尺寸（像素），0 表示不缩放。                                          |
| **PNG**  | **JPEG** | - `width`, `height`：支持调整输出图片尺寸（像素），0 表示不缩放。<br>- `quality`：JPEG 压缩质量 (1-100)，默认 100。 |
| **SVG**  | **PNG**  | - `width`, `height`：支持调整输出图片尺寸（像素），0 表示使用 SVG 定义的尺寸。                               |
| **SVG**  | **JPEG** | - `width`, `height`：支持调整输出图片尺寸（像素），0 表示使用 SVG 定义的尺寸。<br>- `quality`：JPEG 压缩质量 (1-100)，默认 100。 |
| **PNG**  | **SVG**  | - `width`, `height`：指定 SVG 的显示尺寸（像素），0 表示使用原图尺寸。<br>*(注：采用嵌入式转换，将 PNG 嵌入 SVG)*          |
| **JPEG** | **SVG**  | - `width`, `height`：指定 SVG 的显示尺寸（像素），0 表示使用原图尺寸。<br>*(注：采用嵌入式转换，将 JPEG 嵌入 SVG)*         |

## 🏗 架构概览

Ruyi 的核心由以下几个部分组成：

- **Ruyi (Engine)**: 对外暴露的统一入口，负责协调组件工作。
- **Registry**: 注册中心，维护所有已注册的转换器，支持 O(1) 复杂度查找。
- **Converter**: 具体的转换逻辑实现者，每个转换器负责一对特定格式的转换（单例、无状态）。
- **Concept**: 定义数据的语义（如 `JPEG`, `PNG`），结合 `Kind`（如 `File`）确保转换的语义正确性。

## 🤝 贡献指南

欢迎提交 Issue 或 Pull Request 来丰富 Ruyi 的转换能力！请确保新增的 Converter 遵循 `pkg/contract` 中的接口定义，并附带相应的测试用例。

## 📄 许可证

MIT License
