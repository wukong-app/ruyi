package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wukong-app/ruyi"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

// Config 存储命令行配置参数
type Config struct {
	Kind   string
	From   string
	To     string
	In     string
	Out    string
	Params ParamMap
	Help   bool
}

// Handler 定义转换处理逻辑的接口
type Handler interface {
	Handle(ctx context.Context, r contract.Ruyi, cfg *Config) error
}

// main 是命令行工具的入口函数
// 示例命令:
// go run cmd/ruyi/main.go -kind file -from png -to jpeg -in test/testdata/shop.png -out test/testdata/output/shop.jpeg --param width=1024
func main() {
	if err := run(); err != nil {
		// 统一错误处理
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// run 协调整个程序的执行流程：解析参数 -> 初始化 Ruyi -> 获取 Handler -> 执行转换
func run() error {
	// 1. 解析参数
	cfg, err := parseFlags()
	if err != nil {
		// parseFlags 内部已经处理了 Usage 输出，这里直接返回错误
		return err
	}

	// 2. 创建 Ruyi 实例
	r, err := ruyi.New()
	if err != nil {
		return fmt.Errorf("创建 Ruyi 实例失败: %w", err)
	}

	// 3. 获取对应 Kind 的 Handler
	handler, err := getHandler(cfg.Kind)
	if err != nil {
		return err
	}

	// 4. 执行处理逻辑
	return handler.Handle(context.Background(), r, cfg)
}

// parseFlags 解析命令行参数并进行基本的校验
func parseFlags() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.Kind, "kind", "", "转换类型 (file)")
	flag.StringVar(&cfg.From, "from", "", "源 Concept 格式 (例如 png, usd, yyyy-mm-dd)")
	flag.StringVar(&cfg.To, "to", "", "目标 Concept 格式 (例如 jpeg, cny, timestamp)")
	flag.StringVar(&cfg.In, "in", "", "输入内容: 文件路径 或 原始数据")
	flag.StringVar(&cfg.Out, "out", "", "输出内容: 文件路径 或 原始数据输出路径")
	flag.Var(&cfg.Params, "param", "转换器参数（key=value 或 key=value;key=value）可多次指定，或使用分号分隔")
	flag.BoolVar(&cfg.Help, "help", false, "显示帮助信息")

	flag.Parse()

	// 如果指定了 help，只需要必要的参数即可
	if cfg.Help {
		if cfg.Kind == "" || cfg.From == "" || cfg.To == "" {
			fmt.Println("使用 --help 查询参数时，必须提供 kind, from, to 参数")
			fmt.Println("示例: go run cmd/ruyi/main.go -kind file -from png -to jpeg --help")
			flag.Usage()
			return nil, fmt.Errorf("查询参数缺失")
		}
		return cfg, nil
	}

	if cfg.Kind == "" || cfg.From == "" || cfg.To == "" || cfg.In == "" || cfg.Out == "" {
		fmt.Println("必须提供 kind, from, to, in, out 参数")
		flag.Usage()
		return nil, fmt.Errorf("参数缺失")
	}

	return cfg, nil
}

// getHandler 根据 kind 返回对应的 Handler
func getHandler(kind string) (Handler, error) {
	switch strings.ToLower(kind) {
	case "file":
		return &FileHandler{}, nil
	default:
		return nil, fmt.Errorf("未知 kind 类型: %s", kind)
	}
}

// FileHandler 实现文件类型的转换处理逻辑
type FileHandler struct{}

func (h *FileHandler) Handle(ctx context.Context, r contract.Ruyi, cfg *Config) error {
	fromName := contract.ConceptName(cfg.From)
	toName := contract.ConceptName(cfg.To)
	kind := contract.File

	// 获取 Converter
	converter, err := r.GetConverter(ctx, kind, fromName, toName)
	if err != nil {
		if exception.Is(err, exception.ErrNoSupportedConverter) {
			return fmt.Errorf("不支持 %s 类型的转换: %s -> %s", kind, fromName, toName)
		}
		return fmt.Errorf("获取 Converter 失败: %w", err)
	}

	// 如果是 Help 模式，只打印参数信息并退出
	if cfg.Help {
		printParams(converter.Params())
		return nil
	}

	// 打印可用参数信息 (普通模式下也打印，方便调试)
	printParams(converter.Params())

	// 读取输入文件
	fromData, err := os.ReadFile(cfg.In)
	if err != nil {
		return fmt.Errorf("读取输入文件失败: %w", err)
	}

	// 执行转换
	outData, err := converter.Convert(ctx, fromData, cfg.Params)
	if err != nil {
		return fmt.Errorf("文件转换失败: %w", err)
	}

	// 准备输出目录
	dir := filepath.Dir(cfg.Out)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建输出目录(%s)失败: %w", dir, err)
	}

	// 写入输出文件
	if err := os.WriteFile(cfg.Out, outData, 0644); err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	fmt.Printf("转换成功: kind=%s, %s -> %s, 输出: %s\n", kind, fromName, toName, cfg.Out)
	return nil
}

func printParams(params []contract.ConverterParam) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("可用参数 (共 %d 个):\n", len(params)))
	if len(params) == 0 {
		sb.WriteString("  (无参数)\n")
	} else {
		for _, p := range params {
			required := "否"
			if p.Required {
				required = "是"
			}
			// 不显示 Check 函数，因为它无法序列化显示
			checkDesc := ""
			if p.Check != nil {
				checkDesc = " [有值校验]"
			}

			sb.WriteString(fmt.Sprintf(
				"  - %s\n    描述: %s\n    默认值: %q\n    必填: %s%s\n",
				p.Name, p.Desc, p.Default, required, checkDesc,
			))
		}
	}
	fmt.Println(sb.String())
}

// ParamMap 用于解析命令行中的 map 类型参数
type ParamMap map[string]string

func (p *ParamMap) String() string {
	if p == nil || len(*p) == 0 {
		return ""
	}

	var parts []string
	for k, v := range *p {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ", ")
}

func (p *ParamMap) Set(value string) error {
	if *p == nil {
		*p = make(map[string]string)
	}

	// 支持：a=b;c=d
	items := strings.Split(value, ";")

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		kv := strings.SplitN(item, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("参数格式错误，必须是 key=value，得到: %s", item)
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		if key == "" {
			return fmt.Errorf("参数 key 不能为空")
		}

		(*p)[key] = val
	}

	return nil
}
