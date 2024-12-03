package blur

import (
	"image"

	"github.com/disintegration/imaging"
)

type BlurParams struct {
	Sigma float64 `json:"sigma" validate:"required,min=0.1,max=100.0"`
}

func (params *BlurParams) BlurImage(img image.Image) (image.Image, error) {
	return imaging.Blur(img, float64(params.Sigma)), nil
}
