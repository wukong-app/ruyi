package converter

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ contract.Converter = (*pngToJpegConverter)(nil)

// pngToJpegConverter PNG -> JPEG 文件转换器
type pngToJpegConverter struct {
	params contract.ConverterParams
}

func NewPNGToJPEGConverter() contract.Converter {
	params := contract.ConverterParams{}
	params.Append(
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

	return &pngToJpegConverter{
		params: params,
	}
}

func (s *pngToJpegConverter) From() contract.Concept {
	return contract.PNG()
}

func (s *pngToJpegConverter) To() contract.Concept {
	return contract.JPEG()
}

func (s *pngToJpegConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *pngToJpegConverter) Convert(ctx context.Context, in []byte, params map[string]string) ([]byte, error) {
	// 获取参数
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

	// 1、将输入 []byte 转成 io.Reader
	imgReader := bytes.NewReader(in)

	// 2、解码 PNG
	img, err := png.Decode(imgReader)
	if err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "png decode failed")
	}

	// 3、创建新的 RGBA 画布，处理透明背景（PNG 支持透明，JPG 不支持）
	newImg := image.NewRGBA(img.Bounds())

	// 画布涂白色
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// 将 PNG 图像覆盖上去（保留非透明部分）
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)

	var jpegImg image.Image = newImg

	// 缩放
	if width > 0 || height > 0 {
		jpegImg = imaging.Resize(jpegImg, int(width), int(height), imaging.Lanczos)
	}

	// 4、将结果编码为 JPG
	var buf bytes.Buffer
	opt := jpeg.Options{
		Quality: quality,
	}
	if err = jpeg.Encode(&buf, jpegImg, &opt); err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "jpeg encode failed")
	}

	// 5、返回 JPG []byte
	return buf.Bytes(), nil
}
