package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func all_images(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./images")
	if err != nil {
		fmt.Println("Something went wrong", err)
		return
	}
	fmt.Println(files)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Image Gallery</h1>"))

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			thumbnail_path := "/thumb/" + filename
			image_path := "/image/" + filename
			w.Write([]byte(fmt.Sprintf("<a href='%s' target='_blank'><img src='%s'/></a>", image_path, thumbnail_path)))
		}
	}
}

func serve_image(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/image/")
	http.ServeFile(w, r, filepath.Join("./images", filename))
}

func thumbnail(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/thumb/")
	filepath := filepath.Join("./images", filename)
	image, err := load_image(filepath)
	if err != nil {
		fmt.Println("Unable to load image", err)
		http.Error(w, "Unable to load image", http.StatusNotFound)
		return
	}
	thumbnail := resize_image(image, 128, 128)
	w.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(w, thumbnail, nil)
	png.Encode(w, thumbnail)
}

func load_image(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	return image, err
}

func resize_image(img image.Image, height, width int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	// bounds := img.Bounds()
	// w := bounds.Dx()
	// h := bounds.Dy()
	// if (w>=h) {

	// } else {

	// } Will add dynamic thumbnailing later
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * img.Bounds().Max.X / width
			srcY := y * img.Bounds().Max.Y / height
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}
	return resized
}

func main() {
	http.HandleFunc("/", all_images)
	http.HandleFunc("/image/", serve_image)
	http.HandleFunc("/thumb/", thumbnail)

	port := ":8888"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error running server", err)
	} else {
		fmt.Println("Server running on port", port)
	}
}
