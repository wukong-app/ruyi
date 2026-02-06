package register

import (
	"context"
	"fmt"

	"github.com/wukong-app/ruyi/internal/core"
	"github.com/wukong-app/ruyi/pkg/contract"
	"github.com/wukong-app/ruyi/pkg/exception"
)

var _ core.ConverterRegistry = (*converterRegistry)(nil)

func NewConverterRegistry() (core.ConverterRegistry, error) {
	r := &converterRegistry{}
	err := r.Register(converters...)
	if err != nil {
		return nil, exception.Wrapf(err, "Register converters failed")
	}
	return r, nil
}

// converterRegistry 转换器注册器默认实现
type converterRegistry struct {
	// all 所有转换器
	all []core.ConverterAdapter

	// byKind kind -> converter list
	byKind map[contract.Kind][]core.ConverterAdapter

	// matrix kind -> from(Concept name) -> to(Concept name) -> converter
	matrix map[contract.Kind]map[contract.ConceptName]map[contract.ConceptName]core.ConverterAdapter

	// byFrom kind -> from(Concept name) -> converter list
	byFrom map[contract.Kind]map[contract.ConceptName][]core.ConverterAdapter

	// byTo kind -> to(Concept name) -> converter list
	byTo map[contract.Kind]map[contract.ConceptName][]core.ConverterAdapter
}

// Register 注册转换器
// @param converters 转换器列表
func (s *converterRegistry) Register(converters ...core.ConverterAdapter) error {
	if len(converters) == 0 {
		return nil
	}

	var (
		distinct       = make(map[string]struct{})
		genDistinctKey = func(converter core.ConverterAdapter) string {
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

// Find 查找转换器
func (s *converterRegistry) Find(ctx context.Context, kind contract.Kind, from contract.ConceptName, to contract.ConceptName) core.ConverterAdapter {
	// 标准化 name
	fromConcept, ok := core.NormalizeConcept(from)
	if !ok {
		return nil
	}
	toConcept, ok := core.NormalizeConcept(to)
	if !ok {
		return nil
	}
	return s.matrix[kind][fromConcept.Name()][toConcept.Name()]
}

// add 添加转换器
// @receiver s
// @param converter
func (s *converterRegistry) add(converter core.ConverterAdapter) {
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
		s.byKind = make(map[contract.Kind][]core.ConverterAdapter)
	}
	s.byKind[kind] = append(s.byKind[kind], converter)

	// save to matrix
	if s.matrix == nil {
		s.matrix = make(map[contract.Kind]map[contract.ConceptName]map[contract.ConceptName]core.ConverterAdapter)
	}
	if s.matrix[kind] == nil {
		s.matrix[kind] = make(map[contract.ConceptName]map[contract.ConceptName]core.ConverterAdapter)
	}
	if s.matrix[kind][fromName] == nil {
		s.matrix[kind][fromName] = make(map[contract.ConceptName]core.ConverterAdapter)
	}
	s.matrix[kind][fromName][toName] = converter

	// save to byFrom
	if s.byFrom == nil {
		s.byFrom = make(map[contract.Kind]map[contract.ConceptName][]core.ConverterAdapter)
	}
	if s.byFrom[kind] == nil {
		s.byFrom[kind] = make(map[contract.ConceptName][]core.ConverterAdapter)
	}
	s.byFrom[kind][fromName] = append(s.byFrom[kind][fromName], converter)

	// save to byTo
	if s.byTo == nil {
		s.byTo = make(map[contract.Kind]map[contract.ConceptName][]core.ConverterAdapter)
	}
	if s.byTo[kind] == nil {
		s.byTo[kind] = make(map[contract.ConceptName][]core.ConverterAdapter)
	}
	s.byTo[kind][toName] = append(s.byTo[kind][toName], converter)
}
