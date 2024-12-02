package img

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageStorage struct {
	Path string
}

func New(internalStoragePath string) (*ImageStorage, error) {
	const op = "storage.img.New"

	if err := checkFile(internalStoragePath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ImageStorage{Path: internalStoragePath}, nil
}

func (img *ImageStorage) UploadImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	const op = "storage.img.UploadImage"

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

	fileName, err := GenerateName("img", filepath.Ext(handler.Filename))
	if err != nil {
		return "", err
	}

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

func (img *ImageStorage) FindImage(imgName string) (string, error) {
	const op = "storage.img.FindImage"

	filePath := filepath.Join(img.Path, imgName)

	if err := checkFile(img.Path); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return filePath, nil
}

func (img *ImageStorage) DeleteImage(imgName string) error {
	const op = "storage.img.DeleteImage"

	filepath := filepath.Join(img.Path, imgName)

	if err := checkFile(filepath); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (img *ImageStorage) LoadImage(imgName string) (image.Image, error) {
	const op = "storage.img.LoadImage"

	filepath := filepath.Join(img.Path, imgName)

	if err := checkFile(filepath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	loadImg, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return loadImg, nil
}

func (img *ImageStorage) SaveImage(inputImg image.Image, imgName string) (string, error) {
	const op = "storage.img.SaveImage"

	filePath := filepath.Join(img.Path, imgName)

	fileExt := strings.ToLower(filepath.Ext(imgName))
	var err error

	switch fileExt {
	case ".jpg", ".jpeg":
		err = saveJPEG(inputImg, filePath)
	case ".png":
		err = savePNG(inputImg, filePath)
	default:
		return "", fmt.Errorf("%s: unsupported file format: %s", op, fileExt)
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return filePath, nil
}

func GenerateName(prefix string, fileExt string) (string, error) {
	if prefix == "" || fileExt == "" {
		return "", fmt.Errorf("the file prefix or extension must not be empty")
	}

	if !strings.HasPrefix(fileExt, ".") {
		fileExt = "." + fileExt
	}

	return fmt.Sprintf("%s_%s%s", prefix, time.Now().Format("20060102150405"), fileExt), nil
}

func saveJPEG(img image.Image, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}

func savePNG(img image.Image, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func isImage(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}

func checkFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}
	return nil
}
