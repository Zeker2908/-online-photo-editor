package contrast

import (
	"image"

	"github.com/disintegration/imaging"
)

type ContrastParams struct {
	Percentage float64 `json:"percentage" validate:"required,min=-100,max=100"`
}

func (params *ContrastParams) ContrastImage(img image.Image) (image.Image, error) {
	return imaging.AdjustContrast(img, params.Percentage), nil
}
