package save

import (
	"log/slog"
	"mime/multipart"
	"net/http"
	"online-photo-editor/cmd/internal/lib/api/response"
	"online-photo-editor/cmd/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	ImageUrl string `json:"image_url,omitempty"`
}

type ImgSaver interface {
	SaveImg(file multipart.File, handler *multipart.FileHeader) (string, error)
}

func New(log *slog.Logger, imgSaver ImgSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.img.save"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := r.ParseMultipartForm(10 << 20) // 10 MB максимальный размер
		if err != nil {
			log.Error("failed to decode multipart/form-data", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode multipart/form-data"))
			return
		}

		log.Info("multipart/form-data decoded")

		files := r.MultipartForm.File["image"]
		if len(files) != 1 {
			log.Error("invalid number of files uploaded", slog.Int("file_count", len(files)))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("exactly one file must be uploaded"))
			return
		}

		file, handler, err := r.FormFile("image")
		if err != nil {
			log.Error("no file uploaded", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("no file uploaded"))
			return
		}
		defer file.Close()

		imgUrl, err := imgSaver.SaveImg(file, handler)
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

func responseOK(w http.ResponseWriter, r *http.Request, imgUrl string) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{
		Response: response.OK(),
		ImageUrl: imgUrl,
	})
}