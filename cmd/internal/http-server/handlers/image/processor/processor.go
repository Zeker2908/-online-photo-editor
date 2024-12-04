package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"online-photo-editor/cmd/internal/lib/api/blur"
	"online-photo-editor/cmd/internal/lib/api/brightness"
	"online-photo-editor/cmd/internal/lib/api/contrast"
	"online-photo-editor/cmd/internal/lib/api/convert"
	"online-photo-editor/cmd/internal/lib/api/crop"
	"online-photo-editor/cmd/internal/lib/api/gamma"
	"online-photo-editor/cmd/internal/lib/api/resize"
	"online-photo-editor/cmd/internal/lib/api/response"
	"online-photo-editor/cmd/internal/lib/api/saturation"
	"online-photo-editor/cmd/internal/lib/api/sharpen"
	"online-photo-editor/cmd/internal/lib/logger/sl"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const (
	cropAction       = "crop"
	resizeAction     = "resize"
	convertAction    = "convert"
	blurAction       = "blur"
	gammaAction      = "gamma"
	contrastAction   = "contrast"
	sharpenAction    = "sharpen"
	brightnessAction = "brightness"
	saturationAction = "saturation"
)

type ImageAction struct {
	Action string      `json:"action" validate:"required,max=10"`
	Params interface{} `json:"params" validate:"required"`
}

type Request struct {
	Actions   []ImageAction `json:"actions" validate:"required,min=1"`
	ImageName string        `json:"image_name" validate:"required,max=100"`
}

type Response struct {
	response.Response
	ImageUrl string `json:"image_url"`
}

type ImageProcessor interface {
	FindImage(imgName string) (string, error)
	LoadImage(imgName string) (image.Image, error)
	SaveImage(inputImg image.Image, imgName string) (string, error)
	UploadImage(file multipart.File, handler *multipart.FileHeader) (string, error)
	DeleteImage(imgName string) error
	GenerateName(prefix string, fileExt string) (string, error)
}

func New(log *slog.Logger, imgProcessor ImageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.img.processor.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		if !response.Validation(log, w, r, req, http.StatusBadRequest) {
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		imgPath, err := imgProcessor.FindImage(req.ImageName)
		if err != nil {
			log.Error("failed to find image", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("failed to find image"))
			return
		}

		fileExt := strings.ToLower(filepath.Ext(imgPath))

		inputImg, err := imgProcessor.LoadImage(req.ImageName)
		if err != nil {
			log.Error("failed to load image", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("failed to load image"))
			return
		}

		for _, action := range req.Actions {
			if !response.Validation(log, w, r, action, http.StatusBadRequest) {
				return
			}
			switch action.Action {
			case cropAction:
				var params crop.CropParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid crop params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid crop params"))
					return
				}

				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}

				inputImg, err = params.CropImage(inputImg)
			case resizeAction:
				var params resize.ResizeParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid resize params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid resize params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				inputImg, err = params.ResizeImage(inputImg)
			case blurAction:
				var params blur.BlurParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid blur params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid blur params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				inputImg, err = params.BlurImage(inputImg)
			case gammaAction:
				var params gamma.GammaParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid gamma params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid gamma params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				inputImg, err = params.GammaImage(inputImg)
			case contrastAction:
				var params contrast.ContrastParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid gamma params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid gamma params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				inputImg, err = params.ContrastImage(inputImg)
			case sharpenAction:
				var params sharpen.SharpenParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid sharpen params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid sharpen params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				inputImg, err = params.SharpenImage(inputImg)
			case brightnessAction:
				var params brightness.BrightnessParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid sharpen params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid sharpen params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				inputImg, err = params.BrightnessImage(inputImg)
			case saturationAction:
				var params saturation.SaturationParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid sharpen params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid sharpen params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				inputImg, err = params.SaturationImage(inputImg)
			case convertAction:
				var params convert.ConvertParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid convert params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid convert params"))
					return
				}
				if !response.Validation(log, w, r, params, http.StatusBadRequest) {
					return
				}
				fileExt, err = params.ConvertImage()
			default:
				err = fmt.Errorf("field %s must be one of the allowed values`", action.Action)
				log.Error("invalid action", sl.Err(err))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(err.Error()))
				return
			}
			if err != nil {
				log.Error("failed to perform action", sl.Err(err))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(fmt.Sprintf("failed to perform action %s: %v", action.Action, err)))
				return
			}
		}

		imgName, err := imgProcessor.GenerateName("proc", fileExt)
		if err != nil {
			log.Error("failed to generate name", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to generate name"))
			return
		}

		imgUrl, err := imgProcessor.SaveImage(inputImg, imgName)
		if err != nil {
			log.Error("failed to save image", sl.Err(err))
			render.Status(r, http.StatusUnsupportedMediaType)
			render.JSON(w, r, response.Error("failed to save image"))
			return
		}

		log.Info("image saved", slog.String("image url", imgUrl))

		responseOK(w, r, imgUrl)
	}
}

func decodeParams(input interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

func responseOK(w http.ResponseWriter, r *http.Request, imgUrl string) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{
		Response: response.OK(),
		ImageUrl: imgUrl,
	})
}
