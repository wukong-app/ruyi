package converter

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"

	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewGIFToJPEGConverter() contract.Converter {
	return NewBaseConverter(
		contract.GIF(),
		contract.JPEG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return gif.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: ParseQualityParam(params)})
		},
		NewQualityParam(),
	)
}
