package converter

import (
	"bytes"
	"image"
	"image/png"

	"github.com/biessek/golang-ico"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewPNGToICOConverter() contract.Converter {
	return NewBaseConverter(
		contract.PNG(),
		contract.ICO(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return png.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return ico.Encode(w, img)
		},
	)
}
