package converter

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/wukong-app/ruyi/pkg/contract"
	"golang.org/x/image/tiff"
)

func NewTIFFToJPEGConverter() contract.Converter {
	return NewBaseConverter(
		contract.TIFF(),
		contract.JPEG(),
		func(r *bytes.Reader, params map[string]string) (image.Image, error) {
			return tiff.Decode(r)
		},
		func(w *bytes.Buffer, img image.Image, params map[string]string) error {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: ParseQualityParam(params)})
		},
		NewQualityParam(),
	)
}
