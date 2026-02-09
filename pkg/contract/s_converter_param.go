package contract

import "github.com/wukong-app/ruyi/pkg/exception"

type ConverterParams map[string]ConverterParam

func (s ConverterParams) Append(params ...ConverterParam) {
	for _, param := range params {
		s[param.Name] = param
	}
}

func (s ConverterParams) CheckAndGetParams(params map[string]string) (validParams map[string]string, err error) {
	validParams = make(map[string]string, len(s))
	for key, paramDef := range s {
		value := paramDef.Default
		if v, exist := params[key]; exist {
			err = paramDef.Check(v)
			if err != nil {
				return nil, exception.Wrapf(exception.Join(exception.ErrIllegalConverterParam, err), "converter param [%s=%s] check failed", key, v)
			}
			value = v
		}
		validParams[key] = value
	}
	return validParams, nil
}

// ConverterParam 转换器参数规格
type ConverterParam struct {
	Name     string                   // 参数名
	Desc     string                   // 参数描述
	Default  string                   // 默认值
	Required bool                     // 是否必填。true-是，false-否
	Check    func(value string) error // 参数值校验函数
}

func (s ConverterParam) Clone() ConverterParam {
	return ConverterParam{
		Name:     s.Name,
		Desc:     s.Desc,
		Default:  s.Default,
		Required: s.Required,
		Check:    s.Check,
	}
}
