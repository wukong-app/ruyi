package contract

import "context"

// Ruyi 是瑞意（Ruyi）框架的顶层接口。
//
// 说明：
//
//	Ruyi 提供统一的“转换能力”管理，包括文件格式转换、货币转换、时间/日期转换、数字转换等。
//	接口对外只暴露核心转换功能和能力探测功能，保证类型安全。
//	同时提供少量彩蛋方法用于趣味展示，不影响核心业务逻辑。
//
// 接口分为两大部分：
//
//	1、核心功能：能力探测与执行转换。
//	2、彩蛋功能：趣味性接口，例如获取描述、调整尺寸。
type Ruyi interface {
	// -------------------------------
	// 核心功能
	// -------------------------------

	// CanConvert 判断是否支持指定类型的转换能力
	//
	// 参数:
	//   - ctx: 上下文，用于控制超时、取消等
	//   - kind: 转换类型（Kind），例如文件、货币、时间、数字等
	//   - from: 源 Concept 名称
	//   - to: 目标 Concept 名称
	//
	// 返回值:
	//   - bool: 如果存在对应的 Converter 返回 true，否则返回 false
	//
	// 说明:
	//   可在调用 Convert 之前进行能力探测（Capability Check），
	//   避免调用不支持的转换。函数内部通过注册器查找对应的 Converter。
	CanConvert(ctx context.Context, kind Kind, from ConceptName, to ConceptName) bool

	// Convert 通用转换函数（核心功能）
	//
	// 参数:
	//   - ctx: 上下文，用于控制超时、取消等
	//   - kind: 转换类型（Kind），例如文件、货币、时间、数字等
	//   - fromName: 源 Concept 名称
	//   - toName: 目标 Concept 名称
	//   - fromData: 待转换的数据，统一使用 []byte 传递
	//
	// 返回值:
	//   - toData []byte: 转换后的结果
	//   - err error: 转换失败时返回错误，包括以下情况:
	//       1、exception.ErrNoSupportedConverter 找不到对应 Converter
	//       2、exception.ErrConvertFailed Converter 执行出错
	//       3、exception.ErrInvalidConverterOutput Converter 返回值类型不符合预期
	//
	// 说明:
	//   1、通过注册器查找指定 kind、from -> to 的 Converter，如果未找到则返回 ErrNoSupportedConverter。
	//   2、调用 Converter 执行转换，将 fromData 转换为目标类型。
	//   3、返回 Converter 输出，由外层或封装方法（如 ConvertFile、ConvertCurrency 等）做类型断言与安全检查。
	//   4、调用方应根据 kind 对输出进行相应解析，例如：
	//        - 文件类型直接写入文件
	//        - 文本/数字/货币/时间类型转成字符串解析
	Convert(ctx context.Context, kind Kind, fromName ConceptName, toName ConceptName, fromData []byte) (toData []byte, err error)

	// -------------------------------
	// 彩蛋功能（趣味展示，不影响核心逻辑）
	// -------------------------------

	// GetDescription 获取 Ruyi 的描述信息（彩蛋函数）
	//
	// 返回值:
	//   - desc: 当前 Ruyi 的描述文本
	//
	// 说明:
	//   用于趣味展示，与核心转换功能无关。
	GetDescription() (desc string)

	// GetSize 获取 Ruyi 的当前尺寸（彩蛋函数）
	//
	// 返回值:
	//   - size: 当前尺寸值
	//
	// 说明:
	//   彩蛋函数，仅用于趣味展示。
	GetSize() (size int32)

	// Expand 扩展 Ruyi 的尺寸（彩蛋函数）
	//
	// 返回值:
	//   - size: 扩展后的尺寸
	//   - err: 无法继续扩展时返回错误
	//
	// 说明:
	//   使用互斥锁保证线程安全。仅趣味展示，不影响核心转换逻辑。
	Expand() (size int32, err error)

	// Shrink 缩小 Ruyi 的尺寸（彩蛋函数）
	//
	// 返回值:
	//   - size: 缩小后的尺寸
	//   - err: 无法继续缩小时返回错误
	//
	// 说明:
	//   使用互斥锁保证线程安全。仅趣味展示，不影响核心转换逻辑。
	Shrink() (size int32, err error)
}
