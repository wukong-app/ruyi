package converter

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/wukong-app/ruyi/pkg/contract"
	"golang.org/x/image/bmp"
)

func NewBMPToJPEGConverter() contract.Converter {
	return NewBaseConverter(
		contract.BMP(),
		contract.JPEG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return bmp.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: ParseQualityParam(params)})
		},
		NewQualityParam(),
	)
}
