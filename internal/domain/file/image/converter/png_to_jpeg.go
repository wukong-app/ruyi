package converter

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"

	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewPNGToJPEGConverter() contract.Converter {
	return NewBaseConverter(
		contract.PNG(),
		contract.JPEG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return png.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			// PNG -> JPEG 特殊处理：透明背景填充白色
			// 检查是否需要处理透明度
			// 注意：此时的 img 可能是经过 resize 的 NRGBA，或者是原始的 PNG 解码结果

			// 创建新的 RGBA 画布，处理透明背景（PNG 支持透明，JPG 不支持）
			newImg := image.NewRGBA(img.Bounds())

			// 画布涂白色
			draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

			// 将原图覆盖上去（保留非透明部分）
			draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)

			return jpeg.Encode(w, newImg, &jpeg.Options{Quality: ParseQualityParam(params)})
		},
		NewQualityParam(),
	)
}
