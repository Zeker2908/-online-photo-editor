package gamma

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"online-photo-editor/internal/http-server/handlers/image/processor"
	"online-photo-editor/internal/lib/api/gamma"
	"online-photo-editor/internal/lib/api/response"
	"online-photo-editor/internal/lib/logger/sl"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	gamma.GammaParams
	ImageName string `json:"image_name" validate:"required,max=100"`
}

type Response struct {
	response.Response
	ImageUrl string `json:"image_url"`
}

func New(log *slog.Logger, imgGamma processor.ImageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.img.crop.New"

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

		inputImg, err := imgGamma.LoadImage(req.ImageName)
		if err != nil {
			log.Error("failed to load image", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("failed to load image"))
			return
		}

		if !response.Validation(log, w, r, req.GammaParams, http.StatusBadRequest) {
			return
		}

		inputImg, err = req.GammaParams.GammaImage(inputImg)
		if err != nil {
			log.Error("failed to crop image", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("failed to crop image"))
			return
		}

		imgName, err := imgGamma.GenerateName("proc", filepath.Ext(req.ImageName))
		if err != nil {
			log.Error("failed to generate name", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to generate name"))
			return
		}

		imgUrl, err := imgGamma.SaveImage(inputImg, imgName)
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

func responseOK(w http.ResponseWriter, r *http.Request, imgUrl string) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{
		Response: response.OK(),
		ImageUrl: imgUrl,
	})
}
