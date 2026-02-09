package internal

import (
	"github.com/wukong-app/ruyi/internal/domain/file/image/converter"
	"github.com/wukong-app/ruyi/pkg/contract"
)

// ProvideConverters 生成所有转换器列表，供 ConverterRegistry 初始化使用
func ProvideConverters() []contract.Converter {
	return []contract.Converter{
		converter.NewPNGToJPEGConverter(),
	}
}
