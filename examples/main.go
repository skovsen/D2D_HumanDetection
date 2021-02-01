package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// fileContents, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	return nil, err
	// }
	// fi, err := file.Stat()
	// if err != nil {
	// 	return nil, err
	// }

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	file.Close()
	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, "daller")
	if err != nil {
		return nil, err
	}

	part.Write([]byte(encoded))

	// for key, val := range params {
	// 	_ = writer.WriteField(key, val)
	// }
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	//return http.NewRequest("POST", uri, body)
	request, err := http.NewRequest("POST", uri, body)

	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}

func main() {
	// path, _ := os.Getwd()
	// path += "/test.pdf"
	// extraParams := map[string]string{
	// 	"title":       "My Document",
	// 	"author":      "Matt Aimonetti",
	// 	"description": "A document with all the Go programming language secrets",
	// }
	request, err := newfileUploadRequest("http://localhost:8888/upload", "data", "../temp-images/test.png")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		var bodyContent []byte
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		fmt.Println(bodyContent)
	}
}
