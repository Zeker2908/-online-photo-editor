package gamma

import (
	"image"

	"github.com/disintegration/imaging"
)

type GammaParams struct {
	Sigma float64 `json:"sigma" validate:"required,min=0.1,max=100.0"`
}

func (params *GammaParams) GammaImage(img image.Image) (image.Image, error) {
	return imaging.AdjustGamma(img, params.Sigma), nil
}
