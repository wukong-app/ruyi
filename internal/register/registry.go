package register

import (
	"fmt"

	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ core.ConverterRegistry = (*converterRegistry)(nil)

func NewConverterRegistry() (core.ConverterRegistry, error) {
	r := &converterRegistry{}
	err := r.register(converters...)
	if err != nil {
		return nil, exception.Wrapf(err, "register converters failed")
	}
	return r, nil
}

// converterRegistry 转换器注册器默认实现
type converterRegistry struct {
	// all 所有转换器
	all []ConverterAdapter

	// byKind kind -> converter list
	byKind map[core.Kind][]ConverterAdapter

	// matrix kind -> from(Concept name) -> to(Concept name) -> converter
	matrix map[core.Kind]map[string]map[string]ConverterAdapter

	// byFrom kind -> from(Concept name) -> converter list
	byFrom map[core.Kind]map[string][]ConverterAdapter

	// byTo kind -> to(Concept name) -> converter list
	byTo map[core.Kind]map[string][]ConverterAdapter
}

// register 注册转换器
// @param converters 转换器列表
func (s *converterRegistry) register(converters ...ConverterAdapter) error {
	if len(converters) == 0 {
		return nil
	}

	var (
		distinct       = make(map[string]struct{})
		genDistinctKey = func(converter ConverterAdapter) string {
			return fmt.Sprintf("%s_%s", converter.From().Name(), converter.To().Name())
		}
	)

	for _, converter := range converters {
		distinctKey := genDistinctKey(converter)
		if _, exist := distinct[distinctKey]; exist {
			return exception.Errorf("duplicate converter: %s", distinctKey)
		}

		// distinct
		distinct[distinctKey] = struct{}{}

		// save
		s.add(converter)
	}

	return nil
}

func (s *converterRegistry) add(converter ConverterAdapter) {
	var (
		from     = converter.From()
		fromName = from.Name()
		to       = converter.To()
		toName   = to.Name()
		kind     = from.Kind()
	)

	// save to all
	s.all = append(s.all, converter)

	// save to byKind
	if s.byKind == nil {
		s.byKind = make(map[core.Kind][]ConverterAdapter)
	}
	s.byKind[kind] = append(s.byKind[kind], converter)

	// save to matrix
	if s.matrix == nil {
		s.matrix = make(map[core.Kind]map[string]map[string]ConverterAdapter)
	}
	if s.matrix[kind] == nil {
		s.matrix[kind] = make(map[string]map[string]ConverterAdapter)
	}
	if s.matrix[kind][fromName] == nil {
		s.matrix[kind][fromName] = make(map[string]ConverterAdapter)
	}
	s.matrix[kind][fromName][toName] = converter

	// save to byFrom
	if s.byFrom == nil {
		s.byFrom = make(map[core.Kind]map[string][]ConverterAdapter)
	}
	if s.byFrom[kind] == nil {
		s.byFrom[kind] = make(map[string][]ConverterAdapter)
	}
	s.byFrom[kind][fromName] = append(s.byFrom[kind][fromName], converter)

	// save to byTo
	if s.byTo == nil {
		s.byTo = make(map[core.Kind]map[string][]ConverterAdapter)
	}
	if s.byTo[kind] == nil {
		s.byTo[kind] = make(map[string][]ConverterAdapter)
	}
	s.byTo[kind][toName] = append(s.byTo[kind][toName], converter)

}
