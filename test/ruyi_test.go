package ruyi_test

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wukong-app/ruyi"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func TestRuyiExpandAndShrink(t *testing.T) {
	ry, err := ruyi.New()
	if err != nil {
		t.Fatalf("Failed to create Ruyi: %v", err)
	}

	fmt.Printf("Ruyi description is %v \n", ry.GetDescription())
	fmt.Printf("Ruyi size is %v \n", ry.GetSize())

	_, _ = ry.Expand()
	fmt.Printf("Ruyi expanded size is %v \n", ry.GetSize())

	_, _ = ry.Expand()
	fmt.Printf("Ruyi expanded size is %v \n", ry.GetSize())

	_, _ = ry.Shrink()
	fmt.Printf("Ruyi shrunk size is %v \n", ry.GetSize())
}

func TestCanConvert(t *testing.T) {
	t.Run("can convert", func(t *testing.T) {
		ry, err := ruyi.New()
		if err != nil {
			t.Fatalf("Failed to create Ruyi: %v", err)
		}

		result := ry.CanConvert(context.Background(), contract.File, contract.PNG, contract.JPEG)
		require.True(t, result)
	})

	t.Run("cannot convert", func(t *testing.T) {
		ry, err := ruyi.New()
		if err != nil {
			t.Fatalf("Failed to create Ruyi: %v", err)
		}

		result := ry.CanConvert(context.Background(), contract.File, "abc", "123")
		require.False(t, result)
	})

}

func TestConvertFile(t *testing.T) {
	t.Run("PNG 转 JPEG，透明背景转白色", func(t *testing.T) {
		// 1、确保输出目录存在
		if err := os.MkdirAll("testdata/output", 0755); err != nil {
			t.Fatalf("创建输出目录失败: %v", err)
		}

		// 2、读取 PNG 文件内容
		fromData, err := os.ReadFile("testdata/shop.png")
		if err != nil {
			t.Fatalf("读取 PNG 文件失败: %v", err)
		}

		// 3、创建 Ruyi 实例
		ry, err := ruyi.New()
		if err != nil {
			t.Fatalf("创建 Ruyi 实例失败: %v", err)
		}

		ctx := context.Background()
		fromName := contract.PNG
		toName := contract.JPEG

		// 4、检查是否支持转换
		if !ry.CanConvert(ctx, contract.File, fromName, toName) {
			t.Fatalf("不支持 PNG -> JPEG 转换")
		}

		// 5、执行文件转换
		toData, err := ry.ConvertFile(ctx, fromName, toName, fromData)
		if err != nil {
			t.Fatalf("ConvertFile 失败: %v", err)
		}

		// 6、断言：转换结果非空
		if len(toData) == 0 {
			t.Fatalf("转换结果为空")
		}

		// 7、断言：能解码为 JPEG
		img, err := jpeg.Decode(bytes.NewReader(toData))
		if err != nil {
			t.Fatalf("转换结果无法解码为 JPEG: %v", err)
		}
		if img == nil {
			t.Fatalf("解码 JPEG 结果为空")
		}

		// 8、将转换后的 JPG 写入文件（可选）
		outputPath := "testdata/output/shop_with_white_bg.jpg"
		if err := os.WriteFile(outputPath, toData, 0644); err != nil {
			t.Fatalf("写入 JPG 文件失败: %v", err)
		}

		t.Logf("PNG -> JPEG 转换成功，输出文件: %s", outputPath)
	})
}
