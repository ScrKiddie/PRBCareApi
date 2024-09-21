package config

import (
	"github.com/go-playground/validator/v10"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.RegisterValidation("not_contain_space", ValidateNotContainSpace); err != nil {
		log.Fatalln(err)
	}
	if err := v.RegisterValidation("is_password_format", ValidatePasswordFormat); err != nil {
		log.Fatalln(err)
	}
	if err := v.RegisterValidation("image", ValidateImage); err != nil {
		log.Fatalln(err)
	}
	return v
}

func ValidateNotContainSpace(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return !strings.Contains(field, " ")
}

func ValidatePasswordFormat(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>\/?\\|~]`).MatchString(password)
	return hasLower && hasUpper && hasNumber && hasSpecial
}

func ValidateImage(fl validator.FieldLevel) bool {
	file := fl.Field().Interface().(multipart.FileHeader)

	param := fl.Param()
	if param == "" {
		slog.Warn("No validation parameters provided")
		return false
	}

	parts := strings.Split(param, "+")
	if len(parts) != 2 {
		slog.Warn("Invalid validation parameters format", "param", param)
		return false
	}

	dimension := parts[0]
	sizeStr := parts[1]
	dimParts := strings.Split(dimension, "x")
	if len(dimParts) != 2 {
		slog.Warn("Invalid dimension format", "dimension", dimension)
		return false
	}

	width, errW := strconv.Atoi(dimParts[0])
	height, errH := strconv.Atoi(dimParts[1])

	if errW != nil || errH != nil {
		slog.Warn("Invalid dimension values", "width", dimParts[0], "height", dimParts[1], "errorW", errW, "errorH", errH)
		return false
	}

	maxSizeKB, err := strconv.Atoi(sizeStr)
	if err != nil {
		slog.Warn("Invalid max file size value", "maxSizeKB", sizeStr, "error", err)
		return false
	}

	maxFileSize := maxSizeKB * 1024

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		slog.Error("Unsupported file extension", "filename", file.Filename)
		return false
	}

	if file.Size > int64(maxFileSize) {
		slog.Error("File size exceeds limit", "filename", file.Filename, "size", file.Size)
		return false
	}

	imgFile, err := file.Open()
	if err != nil {
		slog.Error("Failed to open file", "filename", file.Filename, "error", err)
		return false
	}
	defer imgFile.Close()

	var img image.Image
	switch ext {
	case ".jpg", ".jpeg", ".png":
		img, _, err = image.Decode(imgFile)
	default:
		slog.Error("Unsupported file format", "filename", file.Filename, "extension", ext)
		return false
	}
	if err != nil {
		slog.Error("Failed to decode image", "filename", file.Filename, "error", err)
		return false
	}

	if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
		slog.Warn("Image dimensions do not match", "filename", file.Filename, "expectedWidth", width, "expectedHeight", height, "actualWidth", img.Bounds().Dx(), "actualHeight", img.Bounds().Dy())
		return false
	}

	return true
}
