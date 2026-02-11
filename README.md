# 如意 (Ruyi)

> 如意金箍棒，随心所欲。Ruyi 能将万物转化为你所需的形态。

Ruyi 是一个基于 Go 语言开发的通用格式转换工具库。它采用插件化架构，旨在为 `wukong` 项目提供强大、灵活且可扩展的数据转换核心。


<p align="center">
<a title="Build Status" target="_blank" href="https://github.com/wukong-app/ruyi/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/wukong-app/ruyi/cd.yml?style=flat-square"></a>
<a title="Releases" target="_blank" href="https://github.com/wukong-app/ruyi/releases"><img src="https://img.shields.io/github/release/wukong-app/ruyi.svg?style=flat-square&color=9CF"></a>
<br>
<a title="Downloads" target="_blank" href="https://github.com/wukong-app/ruyi/releases"><img src="https://img.shields.io/github/downloads/wukong-app/ruyi/total.svg?style=flat-square&color=blueviolet"></a>
<a title="Hits" target="_blank" href="https://github.com/wukong-app/ruyi"><img src="https://hits.b3log.org/wukong-app/ruyi.svg"></a>
<a title="Code Size" target="_blank" href="https://github.com/wukong-app/ruyi"><img src="https://img.shields.io/github/languages/code-size/wukong-app/ruyi.svg?style=flat-square&color=yellow"></a>
<br>
<a title="GitHub Commits" target="_blank" href="https://github.com/wukong-app/ruyi/commits/main"><img src="https://img.shields.io/github/commit-activity/m/wukong-app/ruyi.svg?style=flat-square"></a>
<a title="Last Commit" target="_blank" href="https://github.com/wukong-app/ruyi/commits/main"><img src="https://img.shields.io/github/last-commit/wukong-app/ruyi.svg?style=flat-square&color=FF9900"></a>
<br><br>

## 🌟 核心特性

- **插件化架构**: 转换逻辑封装为独立的 Converter，易于扩展和维护。
- **类型安全**: 基于 Concept（概念）和 Kind（类型）的强类型设计。
- **参数化控制**: 支持在转换过程中传递参数（如图片缩放、质量控制）。
- **统一接口**: 通过统一的 Registry 和 Engine 进行管理，屏蔽底层差异。

## 📦 安装与下载

### 1. 作为 Go 库使用

```bash
go get github.com/wukong-app/ruyi
```

### 2. 下载命令行工具 (CLI)

我们为 **macOS**、**Linux** 和 **Windows** 提供了预编译的二进制文件。

您可以前往 [Releases 页面](https://github.com/wukong-app/ruyi/releases) 下载最新版本。

| 平台 | 架构 | 文件名 |
| :--- | :--- | :--- |
| **macOS** | Intel (amd64) | `ruyi-vX.Y.Z-darwin-amd64` |
| **macOS** | Apple Silicon (arm64) | `ruyi-vX.Y.Z-darwin-arm64` |
| **Linux** | amd64 | `ruyi-vX.Y.Z-linux-amd64` |
| **Windows** | amd64 | `ruyi-vX.Y.Z-windows-amd64.exe` |

*(注: `vX.Y.Z` 为版本号，请下载时替换为实际版本)*

---

## 💻 命令行工具 (CLI) 使用指南

无需编写代码，直接使用命令行工具即可完成转换。

### 🔹 macOS / Linux 使用方法

1.  **下载** 对应系统的二进制文件。
2.  **赋予执行权限** (仅首次)：
    ```bash
    chmod +x ruyi-v1.0.0-darwin-arm64  # 以 macOS arm64 为例
    ```
3.  **运行转换**：
    ```bash
    # 基本格式
    ./ruyi-v1.0.0-darwin-arm64 -kind file -from <src> -to <tgt> -in <input> -out <output>
    
    # 示例：将 PNG 转为 JPEG 并调整大小
    ./ruyi-v1.0.0-darwin-arm64 -kind file -from png -to jpeg \
        -in input.png -out output.jpg \
        --param "width=800;quality=90"
    ```

### 🔹 Windows 使用方法

1.  **下载** `ruyi-vX.Y.Z-windows-amd64.exe`。
2.  打开 **CMD** 或 **PowerShell**，进入文件所在目录。
3.  **运行转换**：
    ```powershell
    # 基本格式
    .\ruyi-vX.Y.Z-windows-amd64.exe -kind file -from <src> -to <tgt> -in <input> -out <output>
    
    # 示例：将 PNG 转为 ICO 图标
    .\ruyi-v1.0.0-windows-amd64.exe -kind file -from png -to ico -in logo.png -out logo.ico
    ```

### ❓ 获取帮助

如果不确定某个转换器支持哪些参数，可以使用 `--help`：

```bash
# 查询 SVG 转 PNG 的可用参数
./ruyi -kind file -from svg -to png --help
```

---

## 🚀 快速开始 (Go SDK)

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

| 源 \ 目标   | PNG | JPEG | GIF | BMP | TIFF | ICO | WEBP | HEIC | SVG |
|:---------|:---:|:----:|:---:|:---:|:----:|:---:|:----:|:----:|:---:|
| **PNG**  |  -  |  ✅   |  ✅  |  ✅  |  ✅   |  ✅  |  ⚠️  |  ⚠️  |  ✅  |
| **JPEG** |  ✅  |  -   |  ✅  |  ✅  |  ✅   |  ✅  |  ⚠️  |  ⚠️  |  ✅  |
| **GIF**  |  ✅  |  ✅   |  -  |  -  |  -   |  -  |  -   |  -   |  -  |
| **BMP**  |  ✅  |  ✅   |  -  |  -  |  -   |  -  |  -   |  -   |  -  |
| **TIFF** |  ✅  |  ✅   |  -  |  -  |  -   |  -  |  -   |  -   |  -  |
| **ICO**  |  ✅  |  ✅   |  -  |  -  |  -   |  -  |  -   |  -   |  -  |
| **WEBP** |  ✅  |  ✅   |  -  |  -  |  -   |  -  |  -   |  -   |  -  |
| **HEIC** |  ✅  |  ✅   |  -  |  -  |  -   |  -  |  -   |  -   |  -  |
| **SVG**  |  ✅  |  ✅   |  -  |  -  |  -   |  -  |  -   |  -   |  -  |

> **注:**
> * ✅: 完全支持
> * ⚠️: 暂不支持

### 🎛️ 通用参数说明

大多数转换器都支持以下通用参数来控制输出结果：

| 参数名           | 说明                           | 适用范围       | 默认值   |
|:--------------|:-----------------------------|:-----------|:------|
| **`width`**   | 输出图片的宽度（像素）。`0` 表示保持原比例或不缩放。 | 所有图片转换     | `0`   |
| **`height`**  | 输出图片的高度（像素）。`0` 表示保持原比例或不缩放。 | 所有图片转换     | `0`   |
| **`quality`** | 图片压缩质量 (1-100)，值越高画质越好，文件越大。 | JPEG, WEBP | `100` |

*提示：使用 CLI 工具时，可以通过 `go run cmd/ruyi/main.go -kind file -from <src> -to <tgt> --help`
查看特定转换器的详细参数。*

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

* **[imaging](https://github.com/disintegration/imaging)**: 提供核心的图片处理算法（缩放、旋转等）。
* **[golang.org/x/image](https://pkg.go.dev/golang.org/x/image)**: 提供 BMP, TIFF, WEBP 等格式的编解码支持。
* **[goheif](https://github.com/jdeng/goheif)**: 提供 HEIC 格式的纯 Go 解码支持。
* **[oksvg](https://github.com/srwiley/oksvg)**: 提供 SVG 格式的解析和渲染支持。
* **[golang-ico](https://github.com/biessek/golang-ico)**: 提供 ICO 格式的编解码支持。
* **[SiYuan](https://github.com/siyuan-note/siyuan)**: 参考了其工程构建思路。

## 💻 开发指南

### 常用命令

* **编译**: `make build`
* **运行测试**: `make test`
* **生成依赖注入代码**: `make wire`
* **格式化代码**: `make fmt`
* **清理构建产物**: `make clean`

## 📄 许可证

MIT License
