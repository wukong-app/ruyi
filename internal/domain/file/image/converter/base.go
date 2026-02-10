package converter

import (
	"bytes"
	"context"
	"image"

	"github.com/disintegration/imaging"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

// DecodeFunc 定义解码函数签名
type DecodeFunc func(r *bytes.Reader, params map[string]string) (image.Image, error)

// EncodeFunc 定义编码函数签名
type EncodeFunc func(w *bytes.Buffer, img image.Image, params map[string]string) error

// BaseConverter 是一个通用的图片转换器实现，封装了常见的 Convert 流程
type BaseConverter struct {
	from       contract.Concept
	to         contract.Concept
	params     contract.ConverterParams
	decodeFunc DecodeFunc
	encodeFunc EncodeFunc
}

// NewBaseConverter 创建一个新的通用转换器
func NewBaseConverter(from, to contract.Concept, decode DecodeFunc, encode EncodeFunc, extraParams ...contract.ConverterParam) *BaseConverter {
	// 默认添加 Width 和 Height 参数
	params := contract.ConverterParams{}
	params.Append(NewWidthParam(), NewHeightParam())

	// 如果是 JPEG 相关的转换（通常 encodeFunc 需要 quality），可以由调用者通过 extraParams 传入 QualityParam
	// 或者我们在这里判断？为了通用性，我们让调用者显式传递 QualityParam 如果他们需要。
	// 但 Width/Height 是几乎所有图片转换都需要的。

	if len(extraParams) > 0 {
		params.Append(extraParams...)
	}

	return &BaseConverter{
		from:       from,
		to:         to,
		params:     params,
		decodeFunc: decode,
		encodeFunc: encode,
	}
}

func (c *BaseConverter) From() contract.Concept {
	return c.from
}

func (c *BaseConverter) To() contract.Concept {
	return c.to
}

func (c *BaseConverter) Params() []contract.ConverterParam {
	// 返回参数副本
	params := make([]contract.ConverterParam, 0, len(c.params))
	for _, param := range c.params {
		params = append(params, param.Clone())
	}
	return params
}

// Convert 执行标准的转换流程：CheckParams -> Decode -> Resize -> Encode
func (c *BaseConverter) Convert(ctx context.Context, in []byte, params map[string]string) ([]byte, error) {
	// 1. 参数校验
	checkedParams, err := c.params.CheckAndGetParams(params)
	if err != nil {
		return nil, err
	}

	// 2. 解析参数
	width, height := ParseResizeParams(checkedParams)

	// 3. 解码
	img, err := c.decodeFunc(bytes.NewReader(in), checkedParams)
	if err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "image decode failed")
	}

	// 4. 缩放 (Resize)
	// 注意：部分格式（如 HEIC）可能返回 YCbCr，如果直接 Encode 为 PNG 可能会有问题。
	// imaging.Resize 返回的是 NRGBA，这很好。
	// 如果不缩放，我们也应该确保图片格式的兼容性（特别是对于 HEIC -> PNG 这种场景）。
	if width > 0 || height > 0 {
		img = imaging.Resize(img, int(width), int(height), imaging.Lanczos)
	} else {
		// 强制转换为 NRGBA/RGBA 以确保最大兼容性 (解决如 HEIC YCbCr -> PNG 的问题)
		// imaging.Clone 会标准化图像格式
		img = imaging.Clone(img)
	}

	// 5. 编码
	var buf bytes.Buffer
	if err := c.encodeFunc(&buf, img, checkedParams); err != nil {
		return nil, exception.Wrapf(exception.Join(exception.ErrConvertFailed, err), "image encode failed")
	}

	return buf.Bytes(), nil
}
