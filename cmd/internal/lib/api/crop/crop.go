package crop

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	"github.com/go-playground/validator/v10"
)

type CropParams struct {
	X      int `json:"x" validate:"required,min=0"`
	Y      int `json:"y" validate:"required,min=0"`
	Width  int `json:"width" validate:"required,min=1"`
	Height int `json:"height" validate:"required,min=1"`
}

func (params *CropParams) validate(img image.Image) error {
	const op = "api.crop.Validate"

	validate := validator.New()

	err := validate.Struct(params)
	if err != nil {
		return fmt.Errorf("%s validation failed: %v", op, err)
	}

	if params.X+params.Width > img.Bounds().Max.X || params.Y+params.Height > img.Bounds().Max.Y {
		return fmt.Errorf("%s crop area exceeds image boundaries", op)
	}

	return nil
}

func (params *CropParams) HandleCrop(img image.Image) (image.Image, error) {
	if err := params.validate(img); err != nil {
		return nil, err
	}

	rect := image.Rect(params.X, params.Y, params.X+params.Width, params.Y+params.Height)
	return imaging.Crop(img, rect), nil
}
