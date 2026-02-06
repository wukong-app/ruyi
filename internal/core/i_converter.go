package core

import "context"

// Converter 泛型接口，定义具体 Concept 间的转换逻辑
type Converter interface {
	// From 返回源 Concept
	From() Concept

	// To 返回目标 Concept
	To() Concept

	// Convert 执行转换
	//
	// 参数:
	//   - ctx: 上下文，用于取消、超时控制
	//   - in: 待转换的值
	//
	// 返回值:
	//   - out: 转换后的值
	//   - err: 转换失败时返回错误
	Convert(ctx context.Context, in []byte) (out []byte, err error)
}
