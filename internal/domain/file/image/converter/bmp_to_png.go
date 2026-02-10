package converter

import (
	"bytes"
	"image"
	"image/png"

	"github.com/wukong-app/ruyi/pkg/contract"
	"golang.org/x/image/bmp"
)

func NewBMPToPNGConverter() contract.Converter {
	return NewBaseConverter(
		contract.BMP(),
		contract.PNG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return bmp.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return png.Encode(w, img)
		},
	)
}
