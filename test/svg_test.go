package ruyi

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wukong-app/ruyi"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func TestSVGConverters(t *testing.T) {
	ry, err := ruyi.New()
	require.NoError(t, err)

	ctx := context.Background()

	// 1. 测试 SVG -> PNG
	t.Run("SVG to PNG", func(t *testing.T) {
		svgContent := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="red" /></svg>`)

		conv, err := ry.GetConverter(ctx, contract.File, contract.Svg, contract.Png)
		require.NoError(t, err)

		out, err := conv.Convert(ctx, svgContent, nil)
		require.NoError(t, err)
		require.NotEmpty(t, out)

		// 验证输出是否为有效 PNG
		img, err := png.Decode(bytes.NewReader(out))
		require.NoError(t, err)
		assert.Equal(t, 100, img.Bounds().Dx())
		assert.Equal(t, 100, img.Bounds().Dy())
	})

	// 2. 测试 SVG -> JPEG
	t.Run("SVG to JPEG", func(t *testing.T) {
		svgContent := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="blue" /></svg>`)

		conv, err := ry.GetConverter(ctx, contract.File, contract.Svg, contract.Jpeg)
		require.NoError(t, err)

		out, err := conv.Convert(ctx, svgContent, nil)
		require.NoError(t, err)
		require.NotEmpty(t, out)

		// 验证输出是否为有效 JPEG
		img, err := jpeg.Decode(bytes.NewReader(out))
		require.NoError(t, err)
		assert.Equal(t, 100, img.Bounds().Dx())
		assert.Equal(t, 100, img.Bounds().Dy())
	})

	// 3. 测试 PNG -> SVG
	t.Run("PNG to SVG", func(t *testing.T) {
		// 创建一个简单的 PNG (50x50)
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)
		var buf bytes.Buffer
		err := png.Encode(&buf, img)
		require.NoError(t, err)
		pngBytes := buf.Bytes()

		conv, err := ry.GetConverter(ctx, contract.File, contract.Png, contract.Svg)
		require.NoError(t, err)

		// Case 1: 默认 (无参数)
		out, err := conv.Convert(ctx, pngBytes, nil)
		require.NoError(t, err)
		svgStr := string(out)
		assert.Contains(t, svgStr, "<svg")
		assert.Contains(t, svgStr, `width="50"`)
		assert.Contains(t, svgStr, `height="50"`)
		assert.Contains(t, svgStr, `viewBox="0 0 50 50"`)
		assert.Contains(t, svgStr, "data:image/png;base64,")

		// Case 2: 指定宽高 (width=100, height=200)
		params := map[string]string{
			"width":  "100",
			"height": "200",
		}
		out, err = conv.Convert(ctx, pngBytes, params)
		require.NoError(t, err)
		svgStr = string(out)

		// 验证 SVG 标签上的宽高被修改
		assert.Contains(t, svgStr, `width="100"`)
		assert.Contains(t, svgStr, `height="200"`)
		// 验证 viewBox 仍保持原图比例
		assert.Contains(t, svgStr, `viewBox="0 0 50 50"`)
	})

	// 4. 测试 JPEG -> SVG
	t.Run("JPEG to SVG", func(t *testing.T) {
		// 创建一个简单的 JPEG (50x50)
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
		var buf bytes.Buffer
		err := jpeg.Encode(&buf, img, nil)
		require.NoError(t, err)
		jpegBytes := buf.Bytes()

		conv, err := ry.GetConverter(ctx, contract.File, contract.Jpeg, contract.Svg)
		require.NoError(t, err)

		// Case 1: 默认 (无参数)
		out, err := conv.Convert(ctx, jpegBytes, nil)
		require.NoError(t, err)
		svgStr := string(out)
		assert.Contains(t, svgStr, "<svg")
		assert.Contains(t, svgStr, `width="50"`)
		assert.Contains(t, svgStr, `height="50"`)
		assert.Contains(t, svgStr, "data:image/jpeg;base64,")

		// Case 2: 指定宽高 (width=150)
		params := map[string]string{
			"width": "150",
		}
		out, err = conv.Convert(ctx, jpegBytes, params)
		require.NoError(t, err)
		svgStr = string(out)

		// 验证 SVG 标签上的宽高被修改
		assert.Contains(t, svgStr, `width="150"`)
		assert.Contains(t, svgStr, `height="50"`) // height 默认为原图高度

		// 验证 viewBox 仍保持原图比例
		expectedViewBox := fmt.Sprintf(`viewBox="0 0 %d %d"`, 50, 50)
		assert.True(t, strings.Contains(svgStr, expectedViewBox))
	})
}
