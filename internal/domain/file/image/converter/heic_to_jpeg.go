package converter

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/jdeng/goheif"
	"github.com/wukong-app/ruyi/pkg/contract"
)

func NewHEICToJPEGConverter() contract.Converter {
	return NewBaseConverter(
		contract.HEIC(),
		contract.JPEG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return goheif.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: ParseQualityParam(params)})
		},
		NewQualityParam(),
	)
}
