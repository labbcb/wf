package client

import (
	"archive/zip"
	"io"
	"os"
)

func zipFiles(zipFile *os.File, files []string) error {
	// new zip writer
	zipWriter := zip.NewWriter(zipFile)

	// iterate over files to zip
	for _, i := range files {
		file, err := os.Open(i)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := addFileToZip(file, zipWriter); err != nil {
			return err
		}
	}

	// close zip writer
	return zipWriter.Close()
}

func addFileToZip(file *os.File, zipWriter *zip.Writer) error {
	// get file info
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// get zip header for file
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// get writer for header
	w2, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	// write file to zip
	_, err = io.Copy(w2, file)
	if err != nil {
		return err
	}
	return nil
}
