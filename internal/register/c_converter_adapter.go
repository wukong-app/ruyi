package register

import (
	"context"
	"fmt"

	"github.com/wukong-app/ruyi/internal/core"
)

// ConverterAdapter Converter 泛型适配器接口
type ConverterAdapter interface {
	From() core.Concept
	To() core.Concept
	Convert(ctx context.Context, in any) (any, error)
}

// AdaptConverter 将 Converter 适配为 ConverterAdapter
func AdaptConverter[T any](c core.Converter[T]) ConverterAdapter {
	return &adapter[T]{c: c}
}

var _ ConverterAdapter = &adapter[any]{}

// adapter Converter 泛型适配器具体实现
type adapter[T any] struct {
	c core.Converter[T]
}

func (a *adapter[T]) From() core.Concept {
	return a.c.From()
}

func (a *adapter[T]) To() core.Concept {
	return a.c.To()
}

func (a *adapter[T]) Kind() core.Kind {
	return a.c.From().Kind()
}

func (a *adapter[T]) Convert(
	ctx context.Context,
	in any,
) (any, error) {
	v, ok := in.(core.Value[T])
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}
	out, err := a.c.Convert(ctx, v)
	if err != nil {
		return nil, err
	}
	return out, nil
}
