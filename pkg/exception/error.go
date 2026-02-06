package exception

import (
	"errors"
	"fmt"
)

// Errorf 创建错误对象
// @param format 错误信息模板
// @param msgArgs 错误信息参数填充
// @return error 错误对象
func Errorf(format string, msgArgs ...any) error {
	return fmt.Errorf(format, msgArgs...)
}

// Wrapf 包装错误
// @param err 原始错误对象
// @param message 上下文信息模板
// @param msgArgs 上下文信息参数填充
// @return error 新的错误对象
func Wrapf(err error, message string, msgArgs ...any) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(message, msgArgs...), err)
}

// Join 将多个错误对象合并成一个错误对象
// @param errs 错误对象列表
// @return error 合并后的错误对象
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Is 判断 err 树中的任意 error 是否与 target 匹配。
// @param err 原始错误对象
// @param target 目标错误对象
// @return bool 是否匹配
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 查找 err 树中与 target 匹配的第一个错误，如果找到，则将 target 设置为该错误值并返回 true。否则，返回 false
// @param err
// @param target
// @return bool
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
