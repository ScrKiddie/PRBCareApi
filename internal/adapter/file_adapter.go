package adapter

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"prb_care_api/internal/model"
	"strconv"
	"strings"
	"time"
)

type FileAdapter struct {
}

func NewFileAdapter() *FileAdapter {
	return &FileAdapter{}
}
func (s *FileAdapter) StoreImageFromBase64(storePath, base64Image string) (*model.File, error) {
	imgData, err := base64.StdEncoding.DecodeString(extractBase64Data(base64Image))
	if err != nil {
		slog.Error("failed to decode base64 image: " + err.Error())
		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		slog.Error("failed to decode image: " + err.Error())
		return nil, err
	}

	basePath, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	storePath = filepath.Join(basePath, storePath)
	if err := os.MkdirAll(storePath, os.ModePerm); err != nil {
		slog.Error("failed to create directory: " + err.Error())
		return nil, err
	}

	uniqueFilename := uuid.New().String() + "." + format
	fullPath := filepath.Join(storePath, uniqueFilename)

	outFile, err := os.Create(fullPath)
	if err != nil {
		slog.Error("failed to create file: " + err.Error())
		return nil, err
	}
	defer func() {
		outFile.Close()
		if err != nil {
			if removeErr := os.Remove(fullPath); removeErr != nil {
				slog.Error("failed to remove file: " + removeErr.Error())
			} else {
				slog.Info("file removed due to error")
			}
		}
	}()

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, img, nil)
	case "png":
		err = png.Encode(outFile, img)
	default:
		err = errors.New("unsupported image format: " + format)
	}

	if err != nil {
		slog.Error("failed to encode image: " + err.Error())
		return nil, err
	}

	return &model.File{Name: uniqueFilename}, nil
}

func extractBase64Data(base64Str string) string {
	if idx := strings.Index(base64Str, ";base64,"); idx != -1 {
		return base64Str[idx+8:]
	}
	return base64Str
}

func (s *FileAdapter) StoreImage(storePath string, f *model.File) (*model.File, error) {
	file, err := f.FileHeader.Open()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer file.Close()

	basePath, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	storePath = filepath.Join(basePath, storePath)

	if err := os.MkdirAll(storePath, os.ModePerm); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	uniqueFilename := uuid.New().String() + filepath.Ext(f.FileHeader.Filename)
	fullPath := filepath.Join(storePath, uniqueFilename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}
	return &model.File{Name: uniqueFilename}, nil
}

func (s *FileAdapter) DeleteFile(storePath string, f *model.File) error {
	basePath, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	storePath = filepath.Join(basePath, storePath)
	filePath := filepath.Join(storePath, f.Name)

	for attempt := 0; attempt < 5; attempt++ {
		err = os.Remove(filePath)
		if err == nil {
			return nil
		}
		slog.Info("attempt " + strconv.Itoa(attempt+1) + ": failed to delete file " + filePath + ": " + err.Error())
		delay := 5 * time.Second * time.Duration(1<<attempt)
		time.Sleep(delay)
	}

	return fmt.Errorf("failed to delete file %s after multiple attempts: %s", filePath, err.Error())
}

func (s *FileAdapter) DeleteFileAsync(storePath string, f *model.File) {
	go func() {
		err := s.DeleteFile(storePath, f)
		if err != nil {
			slog.Error(err.Error())
		}
	}()
}
