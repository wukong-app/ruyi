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

// main 是命令行工具的入口函数
// go run cmd/ruyi/main.go -kind file -from png -to jpeg -in test/testdata/shop.png -out test/testdata/output/shop.jpeg --param width=1024
// go run cmd/ruyi/main.go -kind file -from jpeg -to png -in test/testdata/shop.jpg -out test/testdata/output/shop.png --param width=1024
// go run cmd/ruyi/main.go -kind file -from svg -to png -in test/testdata/shop.svg -out test/testdata/output/shop.png --param width=1024
// go run cmd/ruyi/main.go -kind file -from svg -to jpeg -in test/testdata/shop.svg -out test/testdata/output/shop.jpeg --param width=1024
// go run cmd/ruyi/main.go -kind file -from jpeg -to svg -in test/testdata/shop.jpg -out test/testdata/output/shop.svg --param width=1024
// go run cmd/ruyi/main.go -kind file -from png -to svg -in test/testdata/shop.png -out test/testdata/output/shop.svg --param width=1024
func main() {
	kindFlag := flag.String("kind", "", "转换类型 (file)")
	fromFlag := flag.String("from", "", "源 Concept 格式 (例如 png, usd, yyyy-mm-dd)")
	toFlag := flag.String("to", "", "目标 Concept 格式 (例如 jpeg, cny, timestamp)")
	inFlag := flag.String("in", "", "输入内容: 文件路径 或 原始数据")
	outFlag := flag.String("out", "", "输出内容: 文件路径 或 原始数据输出路径")

	var params ParamMap
	flag.Var(
		&params,
		"param",
		"转换器参数（key=value 或 key=value;key=value）可多次指定，或使用分号分隔",
	)
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

	// 获取 Converter
	converter, err := r.GetConverter(ctx, kind, fromName, toName)
	if err != nil {
		if exception.Is(err, exception.ErrNoSupportedConverter) {
			fmt.Printf("不支持 %s 类型的转换: %s -> %s\n", kind, fromName, toName)
		} else {
			fmt.Printf("获取 Converter 失败: %v\n", err)
		}
		os.Exit(1)
	}

	// 获取 Converter 参数
	converterParams := converter.Params()
	var sb strings.Builder
	sb.WriteString("可用参数:\n")
	for _, p := range converterParams {
		required := "否"
		if p.Required {
			required = "是"
		}
		checkDesc := ""
		if p.Check != nil {
			checkDesc = "（有校验函数）"
		}

		sb.WriteString(fmt.Sprintf(
			"  - %s: %s  默认值: %q  必填: %s %s\n",
			p.Name, p.Desc, p.Default, required, checkDesc,
		))
	}
	fmt.Println(sb.String())

	var outData []byte

	switch kind {
	case contract.File:
		// 文件类型: in 是文件路径
		fromData, err := os.ReadFile(*inFlag)
		if err != nil {
			fmt.Printf("读取输入文件失败: %v\n", err)
			os.Exit(1)
		}

		outData, err = converter.Convert(ctx, fromData, params)
		if err != nil {
			fmt.Printf("文件转换失败: %v\n", err)
			os.Exit(1)
		}

		dir := filepath.Dir(*outFlag)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("创建输出目录(%s)失败: %v\n", dir, err)
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
