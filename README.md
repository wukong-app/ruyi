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

目前 Ruyi 主要支持以下图片格式的转换。我们通过 **源格式 (Source)** 与 **目标格式 (Target)** 的矩阵来展示支持情况及可用参数。

### ✅ 支持矩阵

| 源 \ 目标 | PNG | JPEG | GIF | BMP | TIFF | ICO | WEBP | HEIC | SVG |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |:----:|:---:|
| **PNG** | - | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️|  ⚠️  |  ✅  |
| **JPEG** | ✅ | - | ✅ | ✅ | ✅ | ✅ | ⚠️|  ⚠️  |  ✅  |
| **GIF** | ✅ | ✅ | - | - | - | - | - |  -   |  -  |
| **BMP** | ✅ | ✅ | - | - | - | - | - |  -   |  -  |
| **TIFF** | ✅ | ✅ | - | - | - | - | - |  -   |  -  |
| **ICO** | ✅ | ✅ | - | - | - | - | - |  -   |  -  |
| **WEBP** | ✅ | ✅ | - | - | - | - | - |  -   |  -  |
| **HEIC** | ✅ | ✅ | - | - | - | - | - |  -   |  -  |
| **SVG** | ✅ | ✅ | - | - | - | - | - |  -   |  -  |

> **注:**
> *   ✅: 完全支持
> *   ⚠️: 暂不支持

### 🎛️ 通用参数说明

大多数转换器都支持以下通用参数来控制输出结果：

| 参数名 | 说明 | 适用范围 | 默认值 |
| :--- | :--- | :--- | :--- |
| **`width`** | 输出图片的宽度（像素）。`0` 表示保持原比例或不缩放。 | 所有图片转换 | `0` |
| **`height`** | 输出图片的高度（像素）。`0` 表示保持原比例或不缩放。 | 所有图片转换 | `0` |
| **`quality`** | 图片压缩质量 (1-100)，值越高画质越好，文件越大。 | JPEG, WEBP | `100` |

*提示：使用 CLI 工具时，可以通过 `go run cmd/ruyi/main.go -kind file -from <src> -to <tgt> --help` 查看特定转换器的详细参数。*

## 💻 命令行工具 (CLI)

Ruyi 提供了一个方便的命令行工具，无需编写代码即可直接执行转换。

### 使用方法

```bash
# 1. 运行 CLI 工具 (推荐)
go run cmd/ruyi/main.go -kind file -from <src_format> -to <tgt_format> -in <input_path> -out <output_path> [params...]

# 2. 查询特定转换器的支持参数
go run cmd/ruyi/main.go -kind file -from <src_format> -to <tgt_format> --help
```

### 示例

**1. 将 PNG 转换为 JPEG 并调整尺寸**

```bash
go run cmd/ruyi/main.go -kind file -from png -to jpeg \
    -in test/testdata/shop.png \
    -out output/shop.jpg \
    --param "width=800;quality=90"
```

**2. 查询 SVG 转 PNG 的可用参数**

```bash
go run cmd/ruyi/main.go -kind file -from svg -to png --help
```

---

## 🏗 架构概览

Ruyi 的核心由以下几个部分组成：

- **Ruyi (Engine)**: 对外暴露的统一入口，负责协调组件工作。
- **Registry**: 注册中心，维护所有已注册的转换器，支持 O(1) 复杂度查找。
- **Converter**: 具体的转换逻辑实现者，每个转换器负责一对特定格式的转换（单例、无状态）。
- **Concept**: 定义数据的语义（如 `JPEG`, `PNG`），结合 `Kind`（如 `File`）确保转换的语义正确性。

## 🤝 贡献指南

欢迎提交 Issue 或 Pull Request 来丰富 Ruyi 的转换能力！请确保新增的 Converter 遵循 `pkg/contract` 中的接口定义，并附带相应的测试用例。

## ❤️ 致谢

Ruyi 的强大能力离不开以下优秀的开源项目：

*   **[imaging](https://github.com/disintegration/imaging)**: 提供核心的图片处理算法（缩放、旋转等）。
*   **[golang.org/x/image](https://pkg.go.dev/golang.org/x/image)**: 提供 BMP, TIFF, WEBP 等格式的编解码支持。
*   **[goheif](https://github.com/jdeng/goheif)**: 提供 HEIC 格式的纯 Go 解码支持。
*   **[oksvg](https://github.com/srwiley/oksvg)**: 提供 SVG 格式的解析和渲染支持。
*   **[golang-ico](https://github.com/biessek/golang-ico)**: 提供 ICO 格式的编解码支持。

## 📄 许可证

MIT License
