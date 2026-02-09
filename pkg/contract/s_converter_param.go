package contract

// ConverterParam 转换器参数规格
type ConverterParam struct {
	Name     string // 参数名
	Desc     string // 参数描述
	Default  string // 默认值
	Required bool   // 是否必填。true-是，false-否
}
