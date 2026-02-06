package engine

import (
	"context"
	"math"
	"sync"
	"sync/atomic"

	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ contract.Ruyi = (*Ruyi)(nil)

// NewRuyi 创建并返回一个新的 Ruyi 实例。
//
// 参数:
//   - converterRegistry: Converter 注册中心，用于管理各种转换器
//
// 返回值:
//   - contract.Ruyi 接口类型的实例
//
// 说明:
//
//	NewRuyi 用于对外创建 Ruyi 实例，隐藏内部具体实现。
//	调用方无需关心 Ruyi 的内部细节，只通过接口使用核心功能。
//	如果 converterRegistry 为 nil，则会 panic。
func NewRuyi(converterRegistry core.ConverterRegistry) contract.Ruyi {
	if converterRegistry == nil {
		panic("converterRegistry cannot be nil")
	}

	return &Ruyi{
		description:       "The Ruyi Jingu Bang, Sun Wukong’s magic staff, weighs thirteen thousand five hundred jin.",
		size:              20,
		converterRegister: converterRegistry,
	}
}

// Ruyi 是 contract.Ruyi 接口的具体实现
//
// 说明:
//
//	1、核心功能：提供各种 Concept 间的转换能力。
//	2、彩蛋功能：提供趣味性接口（获取描述、尺寸调整），不影响核心逻辑。
type Ruyi struct {
	//////////////////////////////
	//		Ruyi Jingu Bang（彩蛋属性）	//
	//////////////////////////////
	mx          sync.Mutex // 用于保护 size 读写
	description string     // Ruyi 的描述信息
	size        int32      // 当前尺寸

	//////////////////////////////
	//		依赖组件		//
	//////////////////////////////
	converterRegister core.ConverterRegistry // Converter 注册中心
}

// GetDescription 获取 Ruyi 的描述信息（彩蛋函数）
//
// 返回值:
//   - desc: 当前 Ruyi 的描述文本
//
// 说明:
//
//	此方法用于趣味展示，与核心转换功能无关。
func (s *Ruyi) GetDescription() string {
	return s.description
}

// GetSize 获取 Ruyi 当前尺寸（彩蛋函数）
//
// 返回值:
//   - size: 当前 Ruyi 的尺寸值
//
// 说明:
//
//	仅用于趣味展示，真实业务逻辑不依赖此值。
func (s *Ruyi) GetSize() int32 {
	return atomic.LoadInt32(&s.size)
}

// Expand 扩展 Ruyi 尺寸（彩蛋函数）
//
// 返回值:
//   - size: 扩展后的尺寸
//   - err: 无法继续扩展时返回 ErrRuyiIsBigEnough
//
// 说明:
//
//	使用互斥锁保护 size 线程安全。
//	彩蛋功能，不影响核心转换逻辑。
func (s *Ruyi) Expand() (int32, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.size >= math.MaxInt32 {
		return s.size, exception.ErrRuyiIsBigEnough
	}

	return atomic.AddInt32(&s.size, 1), nil
}

// Shrink 缩小 Ruyi 尺寸（彩蛋函数）
//
// 返回值:
//   - size: 缩小后的尺寸
//   - err: 无法继续缩小时返回 ErrRuyiIsSmallEnough
//
// 说明:
//
//	使用互斥锁保护 size 线程安全。
//	彩蛋功能，仅用于趣味展示。
func (s *Ruyi) Shrink() (int32, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.size <= 0 {
		return s.size, exception.ErrRuyiIsSmallEnough
	}

	return atomic.AddInt32(&s.size, -1), nil
}

// CanConvert 判断是否存在支持指定转换的 Converter
//
// 参数:
//   - ctx: 上下文，用于控制超时、取消等
//   - kind: 转换类型（Kind），例如文件、货币、时间等
//   - from: 源 Concept 名称
//   - to: 目标 Concept 名称
//
// 返回值:
//   - bool: 如果存在对应的 Converter 返回 true，否则返回 false
//
// 说明:
//
//	此函数用于能力探测（Capability Check），可以在调用 Convert 之前
//	先判断是否支持某种转换。函数内部通过注册器查找对应的 Converter。
//	返回 true 表示可以进行转换，返回 false 表示不支持该转换。
func (s *Ruyi) CanConvert(ctx context.Context, kind contract.Kind, from contract.ConceptName, to contract.ConceptName) bool {
	converter := s.converterRegister.Find(ctx, kind, from, to)
	return converter != nil
}

// Convert 通用转换函数（核心功能）
//
// 参数:
//   - ctx: 上下文，用于控制超时、取消等
//   - kind: 转换类型（Kind），例如文件、货币、时间等
//   - fromName: 源 Concept 名称
//   - toName: 目标 Concept 名称
//   - fromData: 待转换的数据（any 类型）
//
// 返回值:
//   - toData: 转换后的数据（any 类型），具体类型由外层或封装方法决定
//   - err: 转换失败时返回的错误，包括找不到转换器或执行失败等
//
// 说明:
//
//	1、首先通过注册器查找指定 kind、from -> to 的 Converter。
//	    如果找不到对应 Converter，则返回 exception.ErrNoSupportedConverter。
//	2、调用 Converter 执行转换，将 fromData 转换为目标类型。
//	    如果转换过程中出现错误，返回 exception.ErrConvertFailed 并封装具体错误信息。
//	3、返回 Converter 输出值，由外层或封装方法（如 ConvertFile、ConvertCurrency）做类型断言和边界安全检查。
//	4、Convert 本身不对输出类型进行断言，以保证通用性。
//	5、外部调用方应根据 kind 或封装方法进行类型安全处理。
func (s *Ruyi) Convert(ctx context.Context, kind contract.Kind, fromName contract.ConceptName, toName contract.ConceptName, fromData any) (toData any, err error) {
	// 1、查找对应 Converter
	converter := s.converterRegister.Find(ctx, kind, fromName, toName)
	if converter == nil {
		return nil, exception.Wrapf(exception.ErrNoSupportedConverter, "converter not found for kind=%s: %s -> %s", kind, fromName, toName)
	}

	// 2、调用 Converter 执行转换
	out, err := converter.Convert(ctx, fromData)
	if err != nil {
		return nil, exception.Wrapf(exception.Join(err, exception.ErrConvertFailed), "conversion failed for kind=%s: %s -> %s", kind, fromName, toName)
	}

	// 3、返回原始输出，由外层或封装方法做类型断言
	return out, nil
}

// ConvertFile 将源 Concept 表示的文件数据转换为目标 Concept 表示的文件数据
//
// 参数:
//   - ctx: 上下文，用于控制超时、取消等
//   - fromName: 源 Concept 名称（例如 JPEG）
//   - toName: 目标 Concept 名称（例如 PNG）
//   - fromData: 源文件的字节数据
//
// 返回值:
//   - []byte: 转换后的目标文件字节数据
//   - error: 转换失败时返回错误，包括以下情况:
//     1、exception.ErrNoSupportedConverter 找不到对应 Converter
//     2、exception.ErrConvertFailed Converter 执行出错
//     3、exception.ErrInvalidConverterOutput Converter 返回值类型不符合预期
//
// 说明:
//
//	1、首先在注册器中查找对应的文件转换器（Converter）。
//	2、调用 Converter 执行转换。
//	3、使用 assertOutput[[]byte] 对 Converter 返回值进行类型断言，
//	    确保边界层类型安全，防止框架契约被破坏。
//
//	注意：
//	  - fromData 和返回值都使用 []byte，适合文件类转换。
//	  - 类型不匹配时会返回异常，标记为内部错误。
func (s *Ruyi) ConvertFile(ctx context.Context, fromName contract.ConceptName, toName contract.ConceptName, fromData []byte) ([]byte, error) {
	// 1、调用通用 Convert 方法
	out, err := s.Convert(ctx, contract.File, fromName, toName, fromData)
	if err != nil {
		return nil, err
	}

	// 2、对输出进行类型断言，保证边界层类型安全
	return assertOutput[[]byte](out, fromName, toName)
}

// assertOutput[T any] 对 Converter 输出进行类型断言，保证边界层类型安全
//
// 参数:
//   - out: Converter 返回值（any 类型）
//   - from: 源 Concept 名称
//   - to: 目标 Concept 名称
//
// 返回值:
//   - T: 断言成功时返回 T 类型的值，失败返回 T 的零值
//   - error: 类型不匹配时返回 ErrInvalidConverterOutput 封装的错误
//
// 说明:
//
//	用于确保 Converter 遵循契约，防止返回不符合预期的类型。
//	错误信息中包含 from/to 信息，方便排查问题。
func assertOutput[T any](out any, from, to contract.ConceptName) (T, error) {
	var zero T
	data, ok := out.(T)
	if ok {
		return data, nil
	}

	return zero, exception.Wrapf(exception.Join(
		exception.ErrInvalidConverterOutput, exception.ErrInternal,
	), "file converter contract violation: %s -> %s, wanted %s,  got %T", from, to, zero, out)
}
