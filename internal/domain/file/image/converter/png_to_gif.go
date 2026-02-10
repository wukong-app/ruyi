package converter

import (
	"bytes"
	"image"
	"image/gif"
	"image/png"

	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewPNGToGIFConverter() contract.Converter {
	return NewBaseConverter(
		contract.PNG(),
		contract.GIF(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return png.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return gif.Encode(w, img, nil)
		},
	)
}
