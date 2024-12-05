package brightness

import (
	"image"

	"github.com/disintegration/imaging"
)

type BrightnessParams struct {
	Percentage float64 `json:"percentage" validate:"required,min=-100,max=100"`
}

func (params *BrightnessParams) BrightnessImage(img image.Image) (image.Image, error) {
	return imaging.AdjustBrightness(img, params.Percentage), nil
}
