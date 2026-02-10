package converter

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ contract.Converter = (*svgToJpegConverter)(nil)

// svgToJpegConverter SVG -> JPEG 文件转换器
type svgToJpegConverter struct {
	params contract.ConverterParams
}

func NewSVGToJPEGConverter() contract.Converter {
	params := contract.ConverterParams{}
	params.Append(NewWidthParam(), NewHeightParam(), NewQualityParam())

	return &svgToJpegConverter{
		params: params,
	}
}

func (s *svgToJpegConverter) From() contract.Concept {
	return contract.SVG()
}

func (s *svgToJpegConverter) To() contract.Concept {
	return contract.JPEG()
}

func (s *svgToJpegConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *svgToJpegConverter) Convert(ctx context.Context, in []byte, params map[string]string) (out []byte, err error) {
	// 1. check params
	params, err = s.params.CheckAndGetParams(params)
	if err != nil {
		return nil, err
	}

	var (
		width   int64
		height  int64
		quality int
	)

	width, _ = strconv.ParseInt(params[core.ParamWidth], 10, strconv.IntSize)
	height, _ = strconv.ParseInt(params[core.ParamHeight], 10, strconv.IntSize)
	quality, _ = strconv.Atoi(params[core.ParamQuality])

	// 2. 解析 SVG
	icon, err := oksvg.ReadIconStream(bytes.NewReader(in))
	if err != nil {
		return nil, exception.Wrapf(err, "svg decode failed")
	}

	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)

	targetW, targetH := w, h
	if width > 0 && height > 0 {
		targetW, targetH = int(width), int(height)
	} else if width > 0 {
		targetW = int(width)
		targetH = int(float64(h) * (float64(width) / float64(w)))
	} else if height > 0 {
		targetH = int(height)
		targetW = int(float64(w) * (float64(height) / float64(h)))
	}

	icon.SetTarget(0, 0, float64(targetW), float64(targetH))

	// 3. 绘制 (RGBA)
	rgba := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	icon.Draw(rasterx.NewDasher(targetW, targetH, rasterx.NewScannerGV(targetW, targetH, rgba, rgba.Bounds())), 1)

	// 4. 处理透明背景 (JPEG 不支持透明，需填充白色背景)
	// 创建一个新的图像，背景填充白色
	bg := image.NewRGBA(rgba.Bounds())
	draw.Draw(bg, bg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	// 将 SVG 绘制结果覆盖上去
	draw.Draw(bg, bg.Bounds(), rgba, rgba.Bounds().Min, draw.Over)

	// 5. 编码为 JPEG
	var buf bytes.Buffer
	// imaging.JPEG 实际上是包装了 jpeg.Encode
	err = imaging.Encode(&buf, bg, imaging.JPEG, imaging.JPEGQuality(quality))
	if err != nil {
		return nil, exception.Wrapf(err, "jpeg encode failed")
	}

	return buf.Bytes(), nil
}
