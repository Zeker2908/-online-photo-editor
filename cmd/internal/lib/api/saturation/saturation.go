package saturation

import (
	"image"

	"github.com/disintegration/imaging"
)

type SaturationParams struct {
	Percentage float64 `json:"percentage" validate:"required,min=-100,max=100"`
}

func (params *SaturationParams) SaturationImage(img image.Image) (image.Image, error) {
	return imaging.AdjustSaturation(img, params.Percentage), nil
}
