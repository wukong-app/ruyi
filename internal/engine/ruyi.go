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
//	NewRuyi 用于对外创建 Ruyi 实例，隐藏内部实现细节。
//	调用方无需关心内部结构，只通过接口使用核心功能。
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

// Ruyi 是 contract.Ruyi 接口的具体实现。
//
// 说明:
//
//	1、核心功能：提供 Concept 间的数据转换能力。
//	2、彩蛋功能：提供趣味性接口（获取描述、尺寸调整），不影响核心逻辑。
type Ruyi struct {
	//////////////////////////////
	//	彩蛋属性
	//////////////////////////////
	mx          sync.Mutex // 用于保护 size 读写
	description string     // Ruyi 的描述信息
	size        int32      // 当前尺寸

	//////////////////////////////
	//	依赖组件
	//////////////////////////////
	converterRegister core.ConverterRegistry // Converter 注册中心
}

func (s *Ruyi) GetConverter(ctx context.Context, kind contract.Kind, from contract.ConceptName, to contract.ConceptName) (contract.Converter, error) {
	// 查找对应 Converter
	converter := s.converterRegister.Find(ctx, kind, from, to)
	if converter == nil {
		return nil, exception.Wrapf(
			exception.ErrNoSupportedConverter,
			"converter not found for kind=%s: %s -> %s",
			kind, from, to,
		)
	}
	return converter, nil
}

// GetDescription 获取 Ruyi 的描述信息（彩蛋函数）
//
// 返回值:
//   - desc: 当前 Ruyi 的描述文本
//
// 说明:
//
//	此方法仅用于趣味展示，不影响核心转换功能。
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
//	此方法仅用于趣味展示，真实业务逻辑不依赖此值。
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

//// CanConvert 判断是否存在支持指定转换的 Converter（核心功能）
////
//// 参数:
////   - ctx: 上下文，用于控制超时、取消等
////   - kind: 转换类型（Kind），例如文件、货币、时间等
////   - from: 源 Concept 名称
////   - to: 目标 Concept 名称
////
//// 返回值:
////   - bool: 如果存在对应的 Converter 返回 true，否则返回 false
////
//// 说明:
////
////	此函数用于能力探测（Capability Check），可在调用 Convert 前判断是否支持某种转换。
//func (s *Ruyi) CanConvert(ctx context.Context, kind contract.Kind, from contract.ConceptName, to contract.ConceptName) bool {
//	converter := s.converterRegister.Find(ctx, kind, from, to)
//	return converter != nil
//}
//
//// Convert 通用转换函数（核心功能）
////
//// 功能说明:
////
////	Convert 是 Ruyi 框架的核心方法，用于执行各种 Concept 间的数据转换。
////	它统一处理不同 Kind 的转换逻辑，例如文件、货币、时间、数字等。
////	函数只负责调用注册的 Converter 并返回原始结果，由外层封装方法负责类型断言与边界安全检查。
////
//// 参数:
////   - ctx: 上下文对象，可用于控制超时、取消等操作。
////   - kind: 转换类型（Kind），例如 File、Currency、Time、Number 等。
////   - fromName: 源 Concept 名称，标识待转换的数据类型。
////   - toName: 目标 Concept 名称，标识转换后的数据类型。
////   - fromData: 待转换的数据，统一使用 []byte 表示。
////
//// 返回值:
////   - toData: 转换后的数据，具体类型由外层封装方法决定。
////   - err: 转换失败时返回的错误，包括找不到转换器或执行失败等。
////
//// 使用说明:
////
////	1、调用前可通过 CanConvert 进行能力探测。
////	2、执行转换时，会通过注册器查找指定 kind、from -> to 的 Converter 并调用其 Convert 方法。
////	3、转换成功后返回 Converter 输出，由外层封装方法进行类型断言。
////	4、不同 Kind 的数据处理方式不同，调用方应使用对应封装方法。
//func (s *Ruyi) Convert(
//	ctx context.Context,
//	kind contract.Kind,
//	fromName contract.ConceptName,
//	toName contract.ConceptName,
//	fromData []byte,
//) (toData []byte, err error) {
//	// 查找对应 Converter
//	converter := s.converterRegister.Find(ctx, kind, fromName, toName)
//	if converter == nil {
//		return nil, exception.Wrapf(
//			exception.ErrNoSupportedConverter,
//			"converter not found for kind=%s: %s -> %s",
//			kind, fromName, toName,
//		)
//	}
//
//	// 调用 Converter 执行转换
//	out, err := converter.Convert(ctx, fromData)
//	if err != nil {
//		return nil, exception.Wrapf(
//			exception.Join(err, exception.ErrConvertFailed),
//			"conversion failed for kind=%s: %s -> %s",
//			kind, fromName, toName,
//		)
//	}
//
//	// 返回原始输出，由外层封装方法做类型断言
//	return out, nil
//}
