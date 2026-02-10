package converter

import (
	"bytes"
	"image"
	"image/png"

	"github.com/wukong-app/ruyi/pkg/contract"
	"golang.org/x/image/webp"
)

func NewWEBPToPNGConverter() contract.Converter {
	return NewBaseConverter(
		contract.WEBP(),
		contract.PNG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return webp.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return png.Encode(w, img)
		},
	)
}
