package converter

import (
	"bytes"
	"image"
	"image/png"

	"github.com/jdeng/goheif"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewHEICToPNGConverter() contract.Converter {
	return NewBaseConverter(
		contract.HEIC(),
		contract.PNG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return goheif.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return png.Encode(w, img)
		},
	)
}
