package converter

import (
	"bytes"
	"context"
	"image"
	"math"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ contract.Converter = (*svgToPngConverter)(nil)

// svgToPngConverter SVG -> PNG 文件转换器
type svgToPngConverter struct {
	params contract.ConverterParams
}

func NewSVGToPNGConverter() contract.Converter {
	params := contract.ConverterParams{}
	params.Append(
		contract.ConverterParam{
			Name:     core.ParamWidth,
			Desc:     "转换后的图片宽度，单位：像素。值为正整数，默认值为 0，表示使用 SVG 定义的宽度。",
			Default:  "0",
			Required: false,
			Check: func(value string) error {
				if value == "" {
					return nil
				}
				v, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					return exception.Wrapf(err, "param value must be a positive integer")
				}
				if v >= math.MaxInt {
					return exception.Errorf("param value must be less than %d", math.MaxInt)
				}
				return nil
			},
		},
		contract.ConverterParam{
			Name:     core.ParamHeight,
			Desc:     "转换后的图片高度，单位：像素。值为正整数，默认值为 0，表示使用 SVG 定义的高度。",
			Default:  "0",
			Required: false,
			Check: func(value string) error {
				if value == "" {
					return nil
				}
				v, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					return exception.Wrapf(err, "param value must be a positive integer")
				}
				if v >= math.MaxInt {
					return exception.Errorf("param value must be less than %d", math.MaxInt)
				}
				return nil
			},
		},
	)

	return &svgToPngConverter{
		params: params,
	}
}

func (s *svgToPngConverter) From() contract.Concept {
	return contract.SVG()
}

func (s *svgToPngConverter) To() contract.Concept {
	return contract.PNG()
}

func (s *svgToPngConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *svgToPngConverter) Convert(ctx context.Context, in []byte, params map[string]string) (out []byte, err error) {
	// 1. check params
	params, err = s.params.CheckAndGetParams(params)
	if err != nil {
		return nil, err
	}

	var (
		width  int64
		height int64
	)

	width, _ = strconv.ParseInt(params[core.ParamWidth], 10, strconv.IntSize)
	height, _ = strconv.ParseInt(params[core.ParamHeight], 10, strconv.IntSize)

	// 2. 解析 SVG
	icon, err := oksvg.ReadIconStream(bytes.NewReader(in))
	if err != nil {
		return nil, exception.Wrapf(err, "svg decode failed")
	}

	// 获取 SVG 原始尺寸
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)

	// 如果指定了宽度或高度，则按比例缩放
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

	// 3. 绘制
	rgba := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	icon.Draw(rasterx.NewDasher(targetW, targetH, rasterx.NewScannerGV(targetW, targetH, rgba, rgba.Bounds())), 1)

	// 4. 编码为 PNG
	var buf bytes.Buffer
	err = imaging.Encode(&buf, rgba, imaging.PNG)
	if err != nil {
		return nil, exception.Wrapf(err, "png encode failed")
	}

	return buf.Bytes(), nil
}
