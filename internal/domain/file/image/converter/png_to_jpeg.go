package converter

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"

	"github.com/wukong-app/ruyi/internal/core"
)

var _ core.Converter = (*pngToJpegConverter)(nil)

// pngToJpegConverter PNG -> JPEG 文件转换器
type pngToJpegConverter struct{}

func NewPNGToJPEGConverter() core.Converter {
	return &pngToJpegConverter{}
}

func (c *pngToJpegConverter) From() core.Concept {
	return core.PNG()
}

func (c *pngToJpegConverter) To() core.Concept {
	return core.JPEG()
}

func (c *pngToJpegConverter) Convert(ctx context.Context, in []byte) ([]byte, error) {
	// 1、将输入 []byte 转成 io.Reader
	imgReader := bytes.NewReader(in)

	// 2、解码 PNG
	img, err := png.Decode(imgReader)
	if err != nil {
		return nil, fmt.Errorf("png decode failed: %w", err)
	}

	// 3、创建新的 RGBA 画布，处理透明背景（PNG 支持透明，JPG 不支持）
	newImg := image.NewRGBA(img.Bounds())

	// 画布涂白色
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// 将 PNG 图像覆盖上去（保留非透明部分）
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)

	// 4、将结果编码为 JPG
	var buf bytes.Buffer
	opt := jpeg.Options{
		Quality: 100,
	}
	if err = jpeg.Encode(&buf, newImg, &opt); err != nil {
		return nil, fmt.Errorf("jpeg encode failed: %w", err)
	}

	// 5、返回 JPG []byte
	return buf.Bytes(), nil
}
