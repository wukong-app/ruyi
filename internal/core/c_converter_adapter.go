package core

import (
	"context"
	"fmt"

	"github.com/wukong-app/ruyi/pkg/contract"
)

// ConverterAdapter Converter 泛型适配器接口
type ConverterAdapter interface {
	From() Concept
	To() Concept
	Kind() contract.Kind
	Convert(ctx context.Context, in any) (any, error)
}

// AdaptConverter 将泛型 Converter[T] 适配为 ConverterAdapter
func AdaptConverter[T any](c Converter[T]) ConverterAdapter {
	return &adapter[T]{c: c}
}

type adapter[T any] struct {
	c Converter[T]
}

func (a *adapter[T]) From() Concept {
	return a.c.From()
}

func (a *adapter[T]) To() Concept {
	return a.c.To()
}

func (a *adapter[T]) Kind() contract.Kind {
	return a.c.From().Kind()
}

func (a *adapter[T]) Convert(ctx context.Context, in any) (any, error) {
	var zero T
	v, ok := in.(T)
	if !ok {
		return nil, fmt.Errorf("converter input type mismatch, expected %T, got %T", zero, in)
	}
	return a.c.Convert(ctx, v)
}
