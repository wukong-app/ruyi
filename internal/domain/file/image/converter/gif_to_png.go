package converter

import (
	"bytes"
	"image"
	"image/gif"
	"image/png"

	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewGIFToPNGConverter() contract.Converter {
	return NewBaseConverter(
		contract.GIF(),
		contract.PNG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return gif.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return png.Encode(w, img)
		},
	)
}
