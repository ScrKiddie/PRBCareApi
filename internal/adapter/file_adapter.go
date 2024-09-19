package adapter

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"prb_care_api/internal/model"
	"strconv"
	"time"
)

type FileAdapter struct {
}

func NewFileAdapter() *FileAdapter {
	return &FileAdapter{}
}

func (s *FileAdapter) StoreImage(storePath string, f *model.FileUpload) (*model.File, error) {
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
	filePath := filepath.Join(storePath, f.Name)

	var err error
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
