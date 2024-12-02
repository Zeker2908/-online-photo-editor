package resize

import (
	"image"

	"github.com/disintegration/imaging"
)

type ResizeParams struct {
	Width  int `json:"width" validate:"required,gte=0"`
	Height int `json:"height" validate:"required,gte=0"`
}

func (params *ResizeParams) Validate()

func (params *ResizeParams) handleResize(img image.Image) (image.Image, error) {
	return imaging.Resize(img, params.Width, params.Height, imaging.Lanczos), nil
}
