package internal

import (
	"github.com/google/wire"
	"github.com/wukong-app/ruyi/internal/domain/file/image/converter"
	"github.com/wukong-app/ruyi/internal/engine"
	"github.com/wukong-app/ruyi/internal/register"
	"github.com/wukong-app/ruyi/pkg/contract"
)

// providerSet combines all dependencies for ruyi
var providerSet = wire.NewSet(
	ProvideConverters,             // 所有 Converter
	register.NewConverterRegistry, // Converter 注册中心
	engine.NewRuyi,                // Ruyi 引擎
)

// ProvideConverters 生成所有转换器列表，供 ConverterRegistry 初始化使用
func ProvideConverters() []contract.Converter {
	return []contract.Converter{
		converter.NewBMPToPNGConverter(),
		converter.NewBMPToJPEGConverter(),
		converter.NewGIFToPNGConverter(),
		converter.NewGIFToJPEGConverter(),
		converter.NewHEICToPNGConverter(),
		converter.NewHEICToJPEGConverter(),
		converter.NewICOToPNGConverter(),
		converter.NewICOToJPEGConverter(),
		converter.NewPNGToGIFConverter(),
		converter.NewPNGToTIFFConverter(),
		//converter.NewPNGToWEBPConverter(),
		converter.NewPNGToICOConverter(),
		//converter.NewPNGToHEICConverter(),
		converter.NewJPEGToPNGConverter(),
		converter.NewJPEGToSVGConverter(),
		converter.NewPNGToJPEGConverter(),
		converter.NewPNGToSVGConverter(),
		converter.NewSVGToPNGConverter(),
		converter.NewSVGToJPEGConverter(),
		converter.NewTIFFToPNGConverter(),
		converter.NewTIFFToJPEGConverter(),
		converter.NewWEBPToPNGConverter(),
		converter.NewWEBPToJPEGConverter(),
	}
}
