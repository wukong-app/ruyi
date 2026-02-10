package converter

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/biessek/golang-ico"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewICOToJPEGConverter() contract.Converter {
	return NewBaseConverter(
		contract.ICO(),
		contract.JPEG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return ico.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: ParseQualityParam(params)})
		},
		NewQualityParam(),
	)
}
