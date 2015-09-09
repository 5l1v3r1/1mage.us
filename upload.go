package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reader, err := r.MultipartReader()
	if err != nil {
		sendUploadError(w, r, "multipart error: "+err.Error())
		return
	}

	ids := []int{}
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err != io.EOF {
				sendUploadError(w, r, "error reading part: "+err.Error())
				return
			}
			break
		}

		if imageId, err := uploadImage(part, r); err != nil {
			sendUploadError(w, r, "error uploading image: "+err.Error())
			return
		} else {
			ids = append(ids, imageId)
		}
	}

	responseMap := map[string]interface{}{"ids": ids}
	if len(ids) == 1 {
		responseMap["identifier"] = ids[0]
	}
	data, _ := json.Marshal(responseMap)
	w.Write(data)
}

func mimeTypeForPart(part *multipart.Part) string {
	// TODO: check the multipart header to see if the browser provides a MIME type
	ext := path.Ext(part.FileName())
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	return mimeType
}

func sendUploadError(w http.ResponseWriter, r *http.Request, errStr string) {
	log.Print("Error from " + r.RemoteAddr + ": " + errStr)
	packet := map[string]string{"error": errStr}
	data, _ := json.Marshal(&packet)
	w.Write(data)
}

func uploadImage(part *multipart.Part, r *http.Request) (id int, err error) {
	if !RateLimitRequest(r) {
		err = errors.New("rate limit exceeded")
		return
	}

	tempFile, err := ioutil.TempFile("", "1mage")
	if err != nil {
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	cappedReader := NewCappedReader(part)
	if _, err := io.Copy(tempFile, cappedReader); err != nil {
		return 0, err
	}

	// TODO: read image file to generate thumbnail.
	// TODO: write thumbnail file.
	// TODO: generate entry in database.

	return 0, nil
}
