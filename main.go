package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"

	"gocv.io/x/gocv"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	log.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	log.Println(r)
	file, _, err := r.FormFile("data")

	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()
	// log.Printf("Uploaded File: %+v\n", handler.Filename)
	// log.Printf("File Size: %+v\n", handler.Size)
	// log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	// read all of the contents of our uploaded file into a
	// byte array
	reader := bufio.NewReader(file)

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	unbased, err := base64.StdEncoding.DecodeString(string(content))
	read := bytes.NewReader(unbased)
	im, err := png.Decode(read)

	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)

	err = png.Encode(buf, im)
	if err != nil {
		panic("Bad png")
	}

	if err != nil {
		panic("Cannot decode b64")
	}

	fileBytes, err := ioutil.ReadAll(buf)
	if err != nil {
		log.Println(err)
	}
	res := recognizeHuman(fileBytes)
	name := "noHIT-*.png"
	if res {
		name = "FOUND-*.png"
	}
	tempFile, err := ioutil.TempFile("temp-images", name)
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!

	if res {
		log.Println("uploaded")
		fmt.Fprintf(w, "Successfully Uploaded File\n")

	} else {
		log.Println("not uploaded")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	dir := http.Dir("html")
	fs := http.FileServer(dir)
	http.Handle("/", fs)
	http.ListenAndServe(":8888", nil)
}

func main() {
	log.Println("Hello World")
	setupRoutes()
}

func recognizeHuman(image []byte) bool {
	// detect faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()
	if !classifier.Load("data/haarcascade_frontalface_default.xml") {
		log.Println("Error reading cascade file: data/haarcascade_frontalface_default.xml")
		return false
	}

	img, err := gocv.IMDecode(image, 1)
	if err != nil {
		log.Fatal(err)
	}

	rects := classifier.DetectMultiScale(img)
	log.Printf("found %d faces\n", len(rects))
	if len(rects) == 0 {
		return false
	}
	// window := gocv.NewWindow("Face Detect")
	// defer window.Close()
	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}
	for _, r := range rects {
		gocv.Rectangle(&img, r, blue, 3)

		// size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		// pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		// gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
	}
	saveFile := "temp-images/found.png"
	gocv.IMWrite(saveFile, img)
	// show the image in the window, and wait 1 millisecond
	// window.IMShow(img)
	// if window.WaitKey(1) >= 0 {
	// 	return true
	// }

	return true

}
