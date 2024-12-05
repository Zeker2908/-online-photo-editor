package convert

type ConvertParams struct {
	Format string `json:"format" validate:"required,lowercase,max=10"`
}

func (params *ConvertParams) ConvertImage() (string, error) {
	return params.Format, nil
}
