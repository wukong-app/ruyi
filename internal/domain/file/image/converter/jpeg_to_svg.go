package converter

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // 注册 jpeg 解码器
	"math"
	"strconv"

	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ contract.Converter = (*jpegToSvgConverter)(nil)

// jpegToSvgConverter JPEG -> SVG 文件转换器 (嵌入式)
type jpegToSvgConverter struct {
	params contract.ConverterParams
}

func NewJPEGToSVGConverter() contract.Converter {
	params := contract.ConverterParams{}
	params.Append(
		contract.ConverterParam{
			Name:     core.ParamWidth,
			Desc:     "SVG 宽度，单位：像素。值为正整数，默认值为 0，表示使用图片原始宽度。",
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
			Desc:     "SVG 高度，单位：像素。值为正整数，默认值为 0，表示使用图片原始高度。",
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

	return &jpegToSvgConverter{
		params: params,
	}
}

func (s *jpegToSvgConverter) From() contract.Concept {
	return contract.JPEG()
}

func (s *jpegToSvgConverter) To() contract.Concept {
	return contract.SVG()
}

func (s *jpegToSvgConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *jpegToSvgConverter) Convert(ctx context.Context, in []byte, params map[string]string) (out []byte, err error) {
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

	// 2. 验证图片有效性并获取尺寸
	img, _, err := image.DecodeConfig(bytes.NewReader(in))
	if err != nil {
		return nil, exception.Wrapf(err, "invalid jpeg image")
	}

	// 3. 确定最终的 SVG 宽高
	targetW := img.Width
	targetH := img.Height

	if width > 0 {
		targetW = int(width)
	}
	if height > 0 {
		targetH = int(height)
	}

	// 4. Base64 编码
	encoded := base64.StdEncoding.EncodeToString(in)

	// 5. 生成 SVG
	// 使用 data URI scheme 嵌入图片
	svgContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="%d" height="%d" viewBox="0 0 %d %d">
<image width="%d" height="%d" xlink:href="data:image/jpeg;base64,%s" />
</svg>`, targetW, targetH, img.Width, img.Height, img.Width, img.Height, encoded)

	return []byte(svgContent), nil
}
