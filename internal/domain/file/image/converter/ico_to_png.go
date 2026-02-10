package converter

import (
	"bytes"
	"image"
	"image/png"

	"github.com/biessek/golang-ico"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewICOToPNGConverter() contract.Converter {
	return NewBaseConverter(
		contract.ICO(),
		contract.PNG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return ico.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return png.Encode(w, img)
		},
	)
}
