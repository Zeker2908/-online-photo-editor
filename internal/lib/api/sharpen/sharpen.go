package sharpen

import (
	"image"

	"github.com/disintegration/imaging"
)

type SharpenParams struct {
	Sigma float64 `json:"sigma" validate:"required,min=0.1,max=100.0"`
}

func (params *SharpenParams) SharpenImage(img image.Image) (image.Image, error) {
	return imaging.Sharpen(img, params.Sigma), nil
}
