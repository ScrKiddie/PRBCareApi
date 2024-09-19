package model

import "mime/multipart"

type FileUpload struct {
	FileHeader *multipart.FileHeader `validate:"image=1200x630+500"`
}
type File struct {
	Name string
}
