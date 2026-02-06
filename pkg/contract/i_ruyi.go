package contract

import "context"

// Ruyi 是瑞意（Ruyi）框架的顶层接口。
//
// 说明：
//
//	Ruyi 接口用于统一管理各种“转换能力”，例如文件格式转换、货币转换、时间转换等。
//	它对外只暴露核心转换功能和能力探测功能，保证类型安全。
//	同时提供少量彩蛋方法用于趣味展示，不影响核心业务逻辑。
//
// 接口分为两大部分：
//
//	1、核心功能：提供能力探测和执行转换的能力。
//	2、彩蛋功能：用于趣味展示，例如获取描述、调整尺寸等。
type Ruyi interface {
	// -------------------------------
	// 核心功能
	// -------------------------------

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
	//   此函数用于能力探测（Capability Check），可以在调用 Convert 之前
	//   先判断是否支持某种转换。函数内部通过注册器查找对应的 Converter。
	//   返回 true 表示可以进行转换，返回 false 表示不支持该转换。
	CanConvert(ctx context.Context, kind Kind, from ConceptName, to ConceptName) bool

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
	Convert(ctx context.Context, kind Kind, fromName ConceptName, toName ConceptName, fromData any) (toData any, err error)

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
	//       1、exception.ErrNoSupportedConverter 找不到对应 Converter
	//       2、exception.ErrConvertFailed Converter 执行出错
	//       3、exception.ErrInvalidConverterOutput Converter 返回值类型不符合预期
	//
	// 说明:
	//   1、首先在注册器中查找对应的文件转换器（Converter）。
	//   2、调用 Converter 执行转换。
	//   3、使用 assertOutput[[]byte] 对 Converter 返回值进行类型断言，
	//       确保边界层类型安全，防止框架契约被破坏。
	//
	//   注意：
	//     - fromData 和返回值都使用 []byte，适合文件类转换。
	//     - 类型不匹配时会返回异常，标记为内部错误。
	ConvertFile(ctx context.Context, fromName ConceptName, toName ConceptName, fromData []byte) (toData []byte, err error)

	// -------------------------------
	// 彩蛋功能（仅趣味展示）
	// -------------------------------

	// GetDescription 获取 Ruyi 的描述信息（彩蛋函数）
	//
	// 返回值:
	//   - desc: 当前 Ruyi 的描述文本
	//
	// 说明:
	//   此方法属于彩蛋函数，仅用于趣味展示，与核心转换功能无关。
	GetDescription() (desc string)

	// GetSize 获取 Ruyi 的当前尺寸（彩蛋函数）
	//
	// 返回值:
	//   - size: 当前 Ruyi 的尺寸值
	//
	// 说明:
	//   彩蛋函数，仅用于趣味展示。
	GetSize() (size int32)

	// Expand 扩展 Ruyi 的尺寸（彩蛋函数）
	//
	// 返回值:
	//   - size: 扩展后的尺寸
	//   - err: 如果无法继续扩展则返回错误
	//
	// 说明:
	//   此方法属于彩蛋功能，与核心业务逻辑无关。
	Expand() (size int32, err error)

	// Shrink 缩小 Ruyi 的尺寸（彩蛋函数）
	//
	// 返回值:
	//   - size: 缩小后的尺寸
	//   - err: 如果无法继续缩小则返回错误
	//
	// 说明:
	//   此方法属于彩蛋功能，与核心业务逻辑无关。
	Shrink() (size int32, err error)
}
