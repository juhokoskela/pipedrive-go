package multipartbody

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// NewFile encodes src as a multipart/form-data body with a single file part.
//
// It returns the Content-Type header value and a *bytes.Reader over the
// encoded body. Because *bytes.Reader is recognized by http.NewRequest,
// the resulting request automatically gets a working GetBody, which
// makes retries (e.g. on 429 / 5xx) safe without any extra wiring.
func NewFile(fieldName, fileName string, src io.Reader) (string, *bytes.Reader, error) {
	return NewFileWithFields(fieldName, fileName, src, nil)
}

// NewFileWithFields encodes src and the supplied form fields as multipart/form-data.
func NewFileWithFields(fieldName, fileName string, src io.Reader, fields url.Values) (string, *bytes.Reader, error) {
	if fieldName == "" {
		return "", nil, fmt.Errorf("multipart field name is required")
	}
	if fileName == "" {
		return "", nil, fmt.Errorf("multipart file name is required")
	}
	if src == nil {
		return "", nil, fmt.Errorf("multipart file content is required")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return "", nil, fmt.Errorf("create multipart file: %w", err)
	}
	if _, err := io.Copy(part, src); err != nil {
		return "", nil, fmt.Errorf("write multipart file: %w", err)
	}
	for name, values := range fields {
		for _, value := range values {
			if err := writer.WriteField(name, value); err != nil {
				return "", nil, fmt.Errorf("write multipart field %q: %w", name, err)
			}
		}
	}
	if err := writer.Close(); err != nil {
		return "", nil, fmt.Errorf("close multipart writer: %w", err)
	}

	return writer.FormDataContentType(), bytes.NewReader(buf.Bytes()), nil
}
