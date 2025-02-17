package rest

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/labstack/echo/v4"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func NewPDFHandler(e *echo.Echo) {
	e.POST("/mergePDF", MergePDF)
	e.POST("/splitPDF", SplitPDF)
	e.POST("/compressPDF", CompressPDF)
}

func MergePDF(c echo.Context) error {
	log.SetPrefix("MergePDF: ")
	// filenames, err := handleInputFile(c)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	// }

	//

	pdfinPath, err := saveFileToTmp(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(pdfinPath) < 2 {
		return c.JSON(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}

	pdfoutPath := "/tmp/mergedPDF.pdf"
	log.Printf("%+v", pdfinPath)
	err = api.MergeCreateFile(pdfinPath, pdfoutPath, false, nil)
	// success return filepath
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// error return e
	return c.Attachment(pdfoutPath, "mergedPDF.pdf")
}

func SplitPDF(c echo.Context) error {
	log.SetPrefix("SplitPDF: ")

	pdfinPath, err := saveFileToTmp(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(pdfinPath) != 1 {
		return c.JSON(http.StatusInternalServerError, domain.ErrInternalServerError.Error())
	}

	pdfoutPath := "/tmp/splitedPDF"
	os.MkdirAll(pdfoutPath, os.ModePerm)
	// log.Printf("%+v", pdfinPath)
	err = api.SplitFile(pdfinPath[0], pdfoutPath, 1, nil)
	if err != nil {
		log.Panicln(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	zipFilePath := "/tmp/splitPDF.zip"
	err = zipFolder(zipFilePath, pdfoutPath)
	if err != nil {
		log.Panicln(err)
		return c.JSON(http.StatusInternalServerError, domain.ErrInternalServerError.Error())
	}
	return c.Attachment(zipFilePath, "splitPDF.zip")
}

func CompressPDF(c echo.Context) error {
	log.SetPrefix("Compress: ")
	pdfinPath, err := saveFileToTmp(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(pdfinPath) != 1 {
		return c.JSON(http.StatusInternalServerError, domain.ErrInternalServerError.Error())
	}

	pdfoutPath := "/tmp/compressedPDF.pdf"
	err = api.OptimizeFile(pdfinPath[0], pdfoutPath, nil)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.Attachment(pdfoutPath, "compressedPDF.pdf")

}

func saveFileToTmp(c echo.Context) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return []string{}, c.String(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}
	// filenames := []string{}
	pdfinPath := []string{}

	files := form.File["files"]
	for index := 0; index < len(files); index++ {
		// filenames = append(filenames, files[index].Filename)

		// file := filenames[index]
		src, err := files[index].Open()
		if err != nil {
			return []string{}, domain.ErrInternalServerError
		}
		filepath := fmt.Sprintf("/tmp/%s", files[index].Filename)
		dst, err := os.Create(filepath)
		if err != nil {
			return []string{}, domain.ErrInternalServerError
		}
		pdfinPath = append(pdfinPath, filepath)
		_, err = io.Copy(dst, src)
		if err != nil {
			return []string{}, domain.ErrInternalServerError
		}
		defer src.Close()
		defer dst.Close()
	}
	return pdfinPath, nil
}

func zipFolder(zipFileName string, inputPath string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	files, err := os.ReadDir(inputPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		filePath := filepath.Join(inputPath, f.Name())

		pageFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer pageFile.Close()

		zipEntry, err := zipWriter.Create(f.Name())
		if err != nil {
			return err
		}

		if _, err := io.Copy(zipEntry, pageFile); err != nil {
			return err
		}
	}

	return nil
}
