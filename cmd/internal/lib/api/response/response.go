package response

import (
	"fmt"
	"log/slog"
	"net/http"
	"online-photo-editor/cmd/internal/lib/logger/sl"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	//TODO: Изменить структуру под изображения

	for _, err := range errs {
		switch err.ActualTag() {
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}

func Validation(log *slog.Logger, w http.ResponseWriter, r *http.Request, s interface{}, errStatus int) bool {
	if err := validator.New().Struct(s); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", sl.Err(err))
		render.Status(r, errStatus)
		render.JSON(w, r, ValidationError(validateErr))

		return false
	}
	return true
}
