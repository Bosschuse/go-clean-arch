package rest_test

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bxcodec/go-clean-arch/internal/rest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func createSomePDF(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString("some context")
	if err != nil {
		return err
	}
	return nil
}

func removeTestPDF(filepath string) error {
	err := os.Remove(filepath)
	return err
}

func createMultipartRequest(url string, fieldName string, filePaths []string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			log.Panicf("failed to open test PDF file: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filePath)
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestMerge(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	err := createSomePDF("file1.pdf")
	if err != nil {
		log.Panicln(err)
	}
	err = createSomePDF("file2.pdf")
	if err != nil {
		log.Panicln(err)
	}

	req, err := createMultipartRequest("/mergePDF", "files", []string{"test.pdf", "test.pdf"})
	if err != nil {
		log.Panicln(err)
	}
	err = removeTestPDF("file1.pdf")
	if err != nil {
		log.Panicln(err)
	}
	err = removeTestPDF("file2.pdf")
	if err != nil {
		log.Panicln(err)
	}
	assert.NoError(t, err)

	c := e.NewContext(req, rec)

	err = rest.MergePDF(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestSpilt(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	// accept 1 pdf file
	req, err := createMultipartRequest("/splitPDF", "files", []string{"test.pdf"})
	if err != nil {
		log.Panicln(err)
	}
	assert.NoError(t, err)
	c := e.NewContext(req, rec)

	err = rest.SplitPDF(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestCompress(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, err := createMultipartRequest("/compressPDF", "files", []string{"test.pdf"})
	if err != nil {
		log.Panicln(err)
	}
	assert.NoError(t, err)
	c := e.NewContext(req, rec)

	err = rest.CompressPDF(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
