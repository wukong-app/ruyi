package contract

import "context"

// Converter 泛型接口，定义具体 Concept 间的转换逻辑
type Converter interface {
	// From 返回源 Concept
	From() Concept

	// To 返回目标 Concept
	To() Concept

	// Convert 通用转换函数（核心功能）
	//
	// 参数:
	//   - ctx: 上下文，用于控制超时、取消等
	//   - in: 待转换的数据，统一使用 []byte 传递
	//
	// 返回值:
	//   - out []byte: 转换后的结果
	//   - err error: 转换失败时返回错误，包括以下情况:
	//       1、exception.ErrConvertFailed Converter 执行出错
	Convert(ctx context.Context, in []byte) (out []byte, err error)
}
