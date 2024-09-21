package model

import "mime/multipart"

type File struct {
	Name       string
	FileHeader *multipart.FileHeader `validate:"image=1200x630+500"`
}
