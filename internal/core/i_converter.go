package core

import "context"

// Converter 转换器
type Converter[T any] interface {

	// From 来源类型
	From() Concept

	// To 目标类型
	To() Concept

	// Convert 执行转换
	Convert(ctx context.Context, in Value[T]) (out Value[T], err error)
}
