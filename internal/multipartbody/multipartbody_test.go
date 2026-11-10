package multipartbody

import (
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"testing"
)

func TestNewFilePreservesFormDataFilenames(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"kuva ä.png", `quote "test".png`} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			contentType, body, err := NewFile("data", name, strings.NewReader("image bytes"))
			if err != nil {
				t.Fatal(err)
			}
			_, params, err := mime.ParseMediaType(contentType)
			if err != nil {
				t.Fatal(err)
			}
			part, err := multipart.NewReader(body, params["boundary"]).NextPart()
			if err != nil {
				t.Fatal(err)
			}
			if got := part.Header.Get("Content-Disposition"); strings.Contains(got, "filename*=") {
				t.Errorf("extended filename parameter is invalid for form-data: %s", got)
			}
			if part.FileName() != name || part.FormName() != "data" {
				t.Fatalf("filename = %q, field = %q", part.FileName(), part.FormName())
			}
			if got := part.Header.Get("Content-Type"); got != "image/png" {
				t.Fatalf("Content-Type = %q", got)
			}
			contents, err := io.ReadAll(part)
			if err != nil || string(contents) != "image bytes" {
				t.Fatalf("body = %q, error = %v", contents, err)
			}
		})
	}
}
