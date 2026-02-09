package contract

import "context"

// Converter 泛型接口，定义具体 Concept 间的转换逻辑
type Converter interface {
	// From 返回源 Concept
	From() Concept

	// To 返回目标 Concept
	To() Concept

	// Params 获取转换器支持的参数列表，
	// 返回转换器参数的副本，返回的切片对于只读使用是安全的。
	Params() []ConverterParam

	// Convert 通用转换函数（核心功能）
	//
	// 参数:
	//   - ctx: 上下文，用于控制超时、取消等
	//   - in: 待转换的数据，统一使用 []byte 传递
	//   - params: 转换参数，由调用方传入，用于控制转换行为。key = 参数名，value = 参数值。参数可选，转换器自身会有默认值
	//
	// 返回值:
	//   - out []byte: 转换后的结果
	//   - err error: 转换失败时返回错误，包括以下情况:
	//       1、exception.ErrConvertFailed Converter 执行出错
	Convert(ctx context.Context, in []byte, params map[string]string) (out []byte, err error)
}
