package converter

import (
	"context"

	"github.com/wukong-app/ruyi/internal/core"
)

var _ core.Converter[[]byte] = (*pngToJpegConverter[[]byte])(nil)

func NewPNGToJPEGConverter() core.Converter[[]byte] {
	return &pngToJpegConverter[[]byte]{}
}

// pngToJpegConverter[T any]  png 转 jpeg
type pngToJpegConverter[T any] struct {
}

func (s *pngToJpegConverter[T]) From() core.Concept {
	return core.PNG()
}

func (s *pngToJpegConverter[T]) To() core.Concept {
	return core.JPEG()
}

func (s *pngToJpegConverter[T]) Convert(ctx context.Context, in core.Value[T]) (out core.Value[T], err error) {
	//TODO implement me
	panic("implement me pngToJpegConverter.Convert")
}
