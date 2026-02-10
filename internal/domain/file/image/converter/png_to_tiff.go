package converter

import (
	"bytes"
	"image"
	"image/png"

	"github.com/wukong-app/ruyi/pkg/contract"
	"golang.org/x/image/tiff"
)

func NewPNGToTIFFConverter() contract.Converter {
	return NewBaseConverter(
		contract.PNG(),
		contract.TIFF(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return png.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return tiff.Encode(w, img, nil)
		},
	)
}
