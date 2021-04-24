package report_portal

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flacatus/che-inspector/pkg/common/clog"
)

const (
	XMLExtensionName = "xml"
)

func (c *API) SendResultsToReportPortal() {
	if !strings.HasSuffix(c.reportPortal.ResultsPath, "/") {
		c.reportPortal.BaseUrl = c.reportPortal.ResultsPath + "/"
	}
	zipFileName := c.reportPortal.ResultsPath + c.reportPortal.Name + ".zip"

	if err := ZipFiles(zipFileName, GetXMLFilesFromDir(c.reportPortal.ResultsPath, XMLExtensionName)); err != nil {
		clog.LOGGER.Fatal("Failed to get junit results from directory %s,%v", c.reportPortal.ResultsPath, err)
	}

	file, err := os.Open(zipFileName)
	if err != nil {
		clog.LOGGER.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	_, _ = io.Copy(part, file)
	_ = writer.Close()

	if _, err = c.Post(context.Background(), "/api/v1/"+c.reportPortal.Project+"/launch/import", writer.FormDataContentType(), body); err != nil {
		clog.LOGGER.Fatal(err)
	}
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func GetXMLFilesFromDir(dir string, ext string) []string {
	var files []string
	filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, dir+"/"+f.Name())
			}
		}
		return nil
	})

	return files
}
