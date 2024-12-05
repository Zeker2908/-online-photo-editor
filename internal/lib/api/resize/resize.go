package resize

import (
	"image"

	"github.com/disintegration/imaging"
)

type ResizeParams struct {
	Width  int `json:"width" validate:"required,min=0,max=8000"`
	Height int `json:"height" validate:"required,min=0,max=8000"`
}

func (params *ResizeParams) ResizeImage(img image.Image) (image.Image, error) {
	return imaging.Resize(img, params.Width, params.Height, imaging.Lanczos), nil
}
