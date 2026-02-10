package converter

import (
	"bytes"
	"image"
	"image/png"

	"github.com/wukong-app/ruyi/pkg/contract"
	"golang.org/x/image/tiff"
)

func NewTIFFToPNGConverter() contract.Converter {
	return NewBaseConverter(
		contract.TIFF(),
		contract.PNG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return tiff.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return png.Encode(w, img)
		},
	)
}
