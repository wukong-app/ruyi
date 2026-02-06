package internal

import (
	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/internal/domain/file/image/converter"
)

// ProvideConverters 生成所有转换器列表，供 ConverterRegistry 初始化使用
func ProvideConverters() []core.Converter {
	return []core.Converter{
		converter.NewPNGToJPEGConverter(),
	}
}
