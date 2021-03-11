package multipartparser

import (
	"errors"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type FormValues struct {
	files  []File
	values map[string]string
}

type File struct {
	header textproto.MIMEHeader
	size   int64
	name   string
	file   multipart.File
}

// ParseMultipart parses the entire multipart request and returns FormValues which contains all files and data.
// You must first extract file before calling this function: r.ParseMultipartForm(10 << 20)
func ParseMultipart(r *http.Request) (FormValues, error) {

	var fvs FormValues

	data := make(map[string]string)

	files := r.MultipartForm.File
	values := r.MultipartForm.Value

	for key := range files {
		file, h, err := r.FormFile(key)
		if err != nil {
			return FormValues{}, errors.New("Error parsing multipart")
		}
		defer file.Close()

		f := File{
			h.Header,
			h.Size,
			h.Filename,
			file,
		}

		fvs.files = append(fvs.files, f)
	}
	for key := range values {
		data[key] = r.FormValue(key)
	}

	fvs.values = data

	return fvs, nil
}

func (fv *FormValues) Files() []File {
	return fv.files
}

func (fv *FormValues) Values() map[string]string {
	return fv.values
}

func (f *File) Header() textproto.MIMEHeader {
	return f.header
}

func (f *File) Size() int64 {
	return f.size
}

func (f *File) Name() string {
	return f.name
}

func (f *File) File() multipart.File {
	return f.file
}
