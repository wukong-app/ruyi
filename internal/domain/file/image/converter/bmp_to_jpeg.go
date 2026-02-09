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
	"golang.org/x/image/bmp"
)

var _ contract.Converter = (*bmpToJpegConverter)(nil)

// bmpToJpegConverter BMP -> JPEG 文件转换器
type bmpToJpegConverter struct {
	params contract.ConverterParams
}

func NewBMPToJPEGConverter() contract.Converter {
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
		contract.ConverterParam{
			Name:     core.ParamQuality,
			Desc:     "将结果编码为 JPG 时的图片质量，范围从 1 到 100（含），越高越好。",
			Default:  "100",
			Required: false,
			Check: func(value string) error {
				if value == "" {
					return exception.Errorf("param is required")
				}
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return exception.Wrapf(err, "param value must be a positive integer")
				}
				if v < 1 || v > 100 {
					return exception.Errorf("param value must be in range [1, 100]")
				}
				return nil
			},
		},
	)

	return &bmpToJpegConverter{
		params: params,
	}
}

func (s *bmpToJpegConverter) From() contract.Concept {
	return contract.BMP()
}

func (s *bmpToJpegConverter) To() contract.Concept {
	return contract.JPEG()
}

func (s *bmpToJpegConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *bmpToJpegConverter) Convert(ctx context.Context, in []byte, params map[string]string) ([]byte, error) {
	params, err := s.params.CheckAndGetParams(params)
	if err != nil {
		return nil, err
	}

	var (
		quality int
		width   int64
		height  int64
	)

	quality, _ = strconv.Atoi(params[core.ParamQuality])
	width, _ = strconv.ParseInt(params[core.ParamWidth], 10, strconv.IntSize)
	height, _ = strconv.ParseInt(params[core.ParamHeight], 10, strconv.IntSize)

	// Decode BMP
	img, err := bmp.Decode(bytes.NewReader(in))
	if err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "bmp decode failed")
	}

	// Resize
	if width > 0 || height > 0 {
		img = imaging.Resize(img, int(width), int(height), imaging.Lanczos)
	}

	// Encode JPEG
	var buf bytes.Buffer
	opt := jpeg.Options{
		Quality: quality,
	}
	if err = jpeg.Encode(&buf, img, &opt); err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "jpeg encode failed")
	}

	return buf.Bytes(), nil
}
