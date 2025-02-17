package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	filenames, err := handleInputFile(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}

	if len(filenames) < 2 {
		return c.JSON(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}

	pdfinPath, err := saveFileToTmp(filenames)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
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
	filenames, err := handleInputFile(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}

	if len(filenames) < 2 {
		return c.JSON(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}

	pdfinPath, err := saveFileToTmp(filenames)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pdfoutPath := "/tmp/splitedPDF.pdf"
	err = api.MergeCreateFile(pdfinPath, pdfoutPath, false, nil)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.Attachment(pdfoutPath, "splitPDF.pdf")
}

func CompressPDF(c echo.Context) error {
	log.SetPrefix("Compress: ")
	filenames, err := handleInputFile(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}

	pdfinPath, err := saveFileToTmp(filenames)
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

func handleInputFile(c echo.Context) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return []string{}, c.String(http.StatusBadRequest, domain.ErrBadParamInput.Error())
	}
	fileList := []string{}
	files := form.File["files"]
	for index := 0; index < len(files); index++ {
		fileList = append(fileList, files[index].Filename)
	}
	return fileList, nil
}

func saveFileToTmp(filenames []string) ([]string, error) {
	pdfinPath := []string{}
	for index := 0; index < len(filenames); index++ {
		file := filenames[index]
		filepath := fmt.Sprintf("/tmp/%s", file)
		dst, err := os.Create(filepath)
		if err != nil {
			return []string{}, domain.ErrInternalServerError
		}
		pdfinPath = append(pdfinPath, filepath)
		defer dst.Close()
	}
	return pdfinPath, nil
}
