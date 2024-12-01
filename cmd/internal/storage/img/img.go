package img

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ImageStorage struct {
	Path string
}

func New(internalStoragePath string) (*ImageStorage, error) {
	const op = "storage.img.New"

	if _, err := os.Stat(internalStoragePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ImageStorage{Path: internalStoragePath}, nil
}

func (img *ImageStorage) SaveImg(file multipart.File, handler *multipart.FileHeader) (string, error) {
	const op = "storage.img.SaveImg"

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	mimeType := http.DetectContentType(buffer)
	if !isImage(mimeType) {
		return "", fmt.Errorf("%s: unsupported file type: %s", op, mimeType)
	}

	fileName := fmt.Sprintf("img_%s%s", time.Now().Format("20060102150405"), filepath.Ext(handler.Filename))
	filePath := filepath.Join(img.Path, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	imageURL := fmt.Sprintf("/images/%s", fileName)

	return imageURL, nil
}

func isImage(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}
