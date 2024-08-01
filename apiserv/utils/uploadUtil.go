package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	customerrors "seclolrip/mchosting/apiserv/errors"
)

func UploadFileToInternal(file *multipart.FileHeader, internalIP string, serverUUID string, filename string) error {
	openedFile, err := file.Open()
	if err != nil {
		return &customerrors.MultipartError{OpenErr: true, Err: err, Source: GetCurrentFuncName()}
	}

	defer openedFile.Close()

	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)

	go func() {
		defer writer.Close()
		defer multipartWriter.Close()

		part, err := multipartWriter.CreateFormFile("file", filename)
		if err != nil {
			writer.CloseWithError(err)
			return
		}

		_, err = io.Copy(part, openedFile)
		if err != nil {
			writer.CloseWithError(err)
			return
		}
	}()

	client := &http.Client{}
	resp, err := client.Post("https://"+internalIP+"/upload/"+serverUUID, multipartWriter.FormDataContentType(), reader)
	if err != nil {
		return &customerrors.HTTPError{StatusCode: uint16(resp.StatusCode), Source: GetCurrentFuncName(), Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return &customerrors.HTTPError{StatusCode: uint16(resp.StatusCode), Source: GetCurrentFuncName(), Err: err}
		}

		return &customerrors.HTTPError{StatusCode: uint16(resp.StatusCode), Source: GetCurrentFuncName(), Err: errors.New(string(bodyBytes))}
	}

	return nil
}
