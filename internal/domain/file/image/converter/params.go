package converter

import (
	"math"
	"strconv"

	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

// CommonParams 定义了图片转换通用的参数名称
const (
	ParamWidth   = core.ParamWidth
	ParamHeight  = core.ParamHeight
	ParamQuality = core.ParamQuality
)

// NewWidthParam 创建宽度参数定义
func NewWidthParam() contract.ConverterParam {
	return contract.ConverterParam{
		Name:     ParamWidth,
		Desc:     "转换后的图片宽度，单位：像素。值为正整数，默认值为 0，表示不缩放。",
		Default:  "0",
		Required: false,
		Check:    CheckPositiveInt,
	}
}

// NewHeightParam 创建高度参数定义
func NewHeightParam() contract.ConverterParam {
	return contract.ConverterParam{
		Name:     ParamHeight,
		Desc:     "转换后的图片高度，单位：像素。值为正整数，默认值为 0，表示不缩放。",
		Default:  "0",
		Required: false,
		Check:    CheckPositiveInt,
	}
}

// NewQualityParam 创建质量参数定义
func NewQualityParam() contract.ConverterParam {
	return contract.ConverterParam{
		Name:     ParamQuality,
		Desc:     "将结果编码为 JPG 时的图片质量，范围从 1 到 100（含），越高越好。",
		Default:  "100",
		Required: false,
		Check:    CheckQuality,
	}
}

// CheckPositiveInt 校验是否为正整数（包含 0）
func CheckPositiveInt(value string) error {
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
}

// CheckQuality 校验图片质量 (1-100)
func CheckQuality(value string) error {
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
}

// ParseResizeParams 解析并返回 width, height 参数
func ParseResizeParams(params map[string]string) (width, height int64) {
	width, _ = strconv.ParseInt(params[ParamWidth], 10, strconv.IntSize)
	height, _ = strconv.ParseInt(params[ParamHeight], 10, strconv.IntSize)
	return
}

// ParseQualityParam 解析并返回 quality 参数
func ParseQualityParam(params map[string]string) int {
	q, _ := strconv.Atoi(params[ParamQuality])
	if q == 0 {
		return 100
	}
	return q
}
