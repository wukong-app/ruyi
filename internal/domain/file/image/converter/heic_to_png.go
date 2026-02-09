package converter

import (
	"bytes"
	"context"
	"image/png"
	"math"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/jdeng/goheif"
	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ contract.Converter = (*heicToPngConverter)(nil)

// heicToPngConverter HEIC -> PNG 文件转换器
type heicToPngConverter struct {
	params contract.ConverterParams
}

func NewHEICToPNGConverter() contract.Converter {
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

	return &heicToPngConverter{
		params: params,
	}
}

func (s *heicToPngConverter) From() contract.Concept {
	return contract.HEIC()
}

func (s *heicToPngConverter) To() contract.Concept {
	return contract.PNG()
}

func (s *heicToPngConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *heicToPngConverter) Convert(ctx context.Context, in []byte, params map[string]string) ([]byte, error) {
	params, err := s.params.CheckAndGetParams(params)
	if err != nil {
		return nil, err
	}

	var (
		width  int64
		height int64
	)

	width, _ = strconv.ParseInt(params[core.ParamWidth], 10, strconv.IntSize)
	height, _ = strconv.ParseInt(params[core.ParamHeight], 10, strconv.IntSize)

	// Decode HEIC
	img, err := goheif.Decode(bytes.NewReader(in))
	if err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "heic decode failed")
	}

	// Resize
	if width > 0 || height > 0 {
		img = imaging.Resize(img, int(width), int(height), imaging.Lanczos)
	}

	// Encode PNG
	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "png encode failed")
	}

	return buf.Bytes(), nil
}
