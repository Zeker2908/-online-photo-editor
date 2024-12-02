package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"
	"net/http"
	"online-photo-editor/cmd/internal/lib/api/crop"
	"online-photo-editor/cmd/internal/lib/api/response"
	"online-photo-editor/cmd/internal/lib/logger/sl"
	imgStorage "online-photo-editor/cmd/internal/storage/img"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	cropAction    = "crop"
	resizeAction  = "resize"
	convertAction = "convert"
)

type ImageAction struct {
	Action string      `json:"action" validate:"required,oneof=crop resize convert"`
	Params interface{} `json:"params" validate:"required"`
}

type Request struct {
	Actions   []ImageAction `json:"actions" validate:"required"`
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
}

func New(log *slog.Logger, imgProcessor ImageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.img.process.New"

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

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))

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

		inputImg, err := imgProcessor.LoadImage(req.ImageName)
		if err != nil {
			log.Error("failed to load image", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("failed to load image"))
			return
		}

		for _, action := range req.Actions {
			switch action.Action {
			case cropAction:
				var params crop.CropParams
				if err := decodeParams(action.Params, &params); err != nil {
					log.Error("invalid crop params", sl.Err(err))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid crop params"))
					return
				}
				inputImg, err = params.HandleCrop(inputImg)
			default:
				log.Error("unknown action", sl.Err(err))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(fmt.Sprintf("unknown action: %s", action.Action)))
				return
			}
			if err != nil {
				log.Error("failed to perform action", sl.Err(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error(fmt.Sprintf("failed to perform action %s: %v", action.Action, err)))
				return
			}
		}

		imgName, err := imgStorage.GenerateName("proc", strings.ToLower(filepath.Ext(imgPath)))
		if err != nil {
			log.Error("failed to generate name", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to generate name"))
			return
		}

		imgUrl, err := imgProcessor.SaveImage(inputImg, imgName)
		if err != nil {
			log.Error("failed to save image", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
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
