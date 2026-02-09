package converter

import (
	"context"
	"math"
	"strconv"

	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
	// Note: We don't have a pure Go WEBP encoder yet, so this converter will fail if implemented naively.
	// However, as per user request to implement converters for png -> webp.
	// Since standard golang.org/x/image/webp only supports Decode, we can't implement Encode without CGO or a 3rd party library.
	// But there is NO mature pure Go WEBP encoder.
	// To satisfy the requirement "compile pass", we can implement the structure but return error or use a placeholder if no library found.
	// BUT, the prompt asked to implement "png -> webp".
	// Let's check if we can find a pure Go webp encoder.
	// Search result said: "WebP ... Encode ❌ No mature pure Go implementation".
	// So we can't really implement a working PNG->WEBP converter without CGO.
	// However, I must follow instructions. I will create the file structure.
	// If I cannot encode, I will return an error "not supported yet" or similar, to allow compilation.
)

var _ contract.Converter = (*pngToWebpConverter)(nil)

// pngToWebpConverter PNG -> WEBP 文件转换器
type pngToWebpConverter struct {
	params contract.ConverterParams
}

func NewPNGToWEBPConverter() contract.Converter {
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

	return &pngToWebpConverter{
		params: params,
	}
}

func (s *pngToWebpConverter) From() contract.Concept {
	return contract.PNG()
}

func (s *pngToWebpConverter) To() contract.Concept {
	return contract.WEBP()
}

func (s *pngToWebpConverter) Params() []contract.ConverterParam {
	params := make([]contract.ConverterParam, 0, len(s.params))
	for _, param := range s.params {
		params = append(params, param.Clone())
	}
	return params
}

func (s *pngToWebpConverter) Convert(ctx context.Context, in []byte, params map[string]string) ([]byte, error) {
	//params, err := s.params.CheckAndGetParams(params)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var (
	//	width  int64
	//	height int64
	//)
	//
	//width, _ = strconv.ParseInt(params[core.ParamWidth], 10, strconv.IntSize)
	//height, _ = strconv.ParseInt(params[core.ParamHeight], 10, strconv.IntSize)
	//
	//// Decode PNG
	//img, err := png.Decode(bytes.NewReader(in))
	//if err != nil {
	//	return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "png decode failed")
	//}
	//
	//// Resize
	//if width > 0 || height > 0 {
	//	img = imaging.Resize(img, int(width), int(height), imaging.Lanczos)
	//}

	// Encode WEBP
	// 由于 Go 生态目前缺乏成熟的 Pure Go WEBP 编码库，
	// 为了保证项目无需 CGO 即可编译通过，这里暂时返回错误。
	// 实际生产中建议使用 CGO 调用 libwebp (如 github.com/chai2010/webp)。
	return nil, exception.Errorf("webp encoding is not supported in pure go mode yet")
}
