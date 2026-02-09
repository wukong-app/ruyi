package core

import (
	"context"

	"github.com/wukong-app/ruyi/pkg/contract"
)

// ConverterRegistry 转换器注册器
type ConverterRegistry interface {

	// Register 注册转换器
	// @param converters 转换器列表
	Register(converters ...contract.Converter) error

	// Find 查找转换器
	// @param kind 种类
	// @param from from名称
	// @param to to名称
	Find(ctx context.Context, kind contract.Kind, from contract.ConceptName, to contract.ConceptName) contract.Converter
}
