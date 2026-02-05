package core

// Value 是泛型包装器，用于统一 Converter 输入输出
type Value[T any] struct {
	Data T
}
