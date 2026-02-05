package register

import (
	"github.com/wukong-app/ruyi/internal/domain/file/image/converter"
)

// converters 转换器列表，供 ConverterRegistry 初始化使用
var converters = []ConverterAdapter{
	AdaptConverter(converter.NewPNGToJPEGConverter()),
}
