package converter

import (
	"bytes"
	"context"
	"image/jpeg"
	"math"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ contract.Converter = (*jpegToPngConverter)(nil)

// jpegToPngConverter JPEG -> PNG 文件转换器
type jpegToPngConverter struct {
	params contract.ConverterParams
}

func NewJPEGToPNGConverter() contract.Converter {
	params := contract.ConverterParams{}
	params.Append(
		contract.ConverterParam{
			Name:     core.ParamWidth,
			Desc:     "转换后的图片宽度，单位：像素。值为正整数，默认值为 0，表示不缩放。",
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
			Desc:     "转换后的图片高度，单位：像素。值为正整数，默认值为 0，表示不缩放。",
			Default:  "0",
			Required: false,
			Check: func(value string) error {
				if value == "" {
					return nil
				}
				v, err := strconv.ParseUint(value, 10, 32)
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

	return &jpegToPngConverter{
		params: params,
	}
}

func (s *jpegToPngConverter) From() contract.Concept {
	return contract.JPEG()
}

func (s *jpegToPngConverter) To() contract.Concept {
	return contract.PNG()
}

func (s *jpegToPngConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *jpegToPngConverter) Convert(ctx context.Context, in []byte, params map[string]string) (out []byte, err error) {
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

	// 2. 解码 (支持多种格式自动识别)
	src, err := jpeg.Decode(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	// 3. 缩放
	if width > 0 || height > 0 {
		src = imaging.Resize(src, int(width), int(height), imaging.Lanczos)
	}

	// 5. 格式转换：编码为 PNG 字节流
	var buf bytes.Buffer
	// imaging.PNG 实际上是包装了标准库的 png.Encode
	err = imaging.Encode(&buf, src, imaging.PNG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
