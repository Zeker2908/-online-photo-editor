package processor_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"net/http/httptest"
	"online-photo-editor/internal/http-server/handlers/image/processor"
	"online-photo-editor/internal/http-server/handlers/image/processor/mocks"
	"online-photo-editor/internal/lib/logger/handlers/slogdiscard"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ProcessImage_Success(t *testing.T) {
	mockProcessor := new(mocks.ImageProcessor)
	logger := slogdiscard.NewDiscardLogger()
	handler := processor.New(logger, mockProcessor)

	reqBody := processor.Request{
		Actions: []processor.ImageAction{
			{Action: "resize", Params: map[string]interface{}{"width": 100, "height": 100}},
		},
		ImageName: "test-image.png",
	}

	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	mockProcessor.On("FindImage", "test-image.png").Return("/path/to/test-image.png", nil)
	mockProcessor.On("LoadImage", "test-image.png").Return(image.NewRGBA(image.Rect(0, 0, 100, 100)), nil)
	mockProcessor.On("GenerateName", "proc", ".png").Return("new-image.png", nil)
	mockProcessor.On("SaveImage", mock.Anything, "new-image.png").Return("/path/to/new-image.png", nil)

	req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response processor.Response
	err = render.DecodeJSON(resp.Body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/new-image.png", response.ImageUrl)
}

func TestHandler_ProcessImage_ImageNotFound(t *testing.T) {
	mockProcessor := new(mocks.ImageProcessor)
	logger := slogdiscard.NewDiscardLogger()
	handler := processor.New(logger, mockProcessor)

	reqBody := processor.Request{
		Actions: []processor.ImageAction{
			{Action: "resize", Params: map[string]interface{}{"width": 100, "height": 100}},
		},
		ImageName: "nonexistent-image.png",
	}

	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	mockProcessor.On("FindImage", "nonexistent-image.png").Return("", errors.New("image not found"))

	req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]string
	err = render.DecodeJSON(resp.Body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "failed to find image", response["error"])
}
