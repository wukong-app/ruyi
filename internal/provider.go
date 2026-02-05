package internal

import (
	"github.com/google/wire"
	"github.com/wukong-app/ruyi/internal/engine"
	"github.com/wukong-app/ruyi/internal/register"
)

// providerSet combines all dependencies for ruyi
var providerSet = wire.NewSet(
	register.NewConverterRegistry, // Converter 注册中心

	// 在此之前添加依赖项↑↑↑↑↑↑↑↑↑↑
	engine.NewRuyi, // Ruyi 引擎。 be sure to put it last！！
	// 在此之后不要添加任何依赖项！！！
)
