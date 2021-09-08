package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var palette = []color.Color{
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.Black,
	color.RGBA{0xff, 0x00, 0x00, 0xff},
}

const (
	greenIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
	redIndex   = 2
)

func main() {
	http.HandleFunc("/", lissajous) // each request calls handler
	http.HandleFunc("/favicon.ico", doNothing)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(w http.ResponseWriter, req *http.Request) {
	var (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	queryParams := req.URL.Query()
	var err error

	if val, ok := queryParams["cycles"]; ok {
		cycles, err = strconv.Atoi(val[0])
	}
	if val, ok := queryParams["res"]; ok {
		res, err = strconv.ParseFloat(val[0], 64)
	}
	if val, ok := queryParams["size"]; ok {
		size, err = strconv.Atoi(val[0])
	}
	if val, ok := queryParams["nframes"]; ok {
		nframes, err = strconv.Atoi(val[0])
	}
	if val, ok := queryParams["delay"]; ok {
		delay, err = strconv.Atoi(val[0])
	}

	if err != nil {
		fmt.Println("error: %v", err)
	}

	fmt.Printf("value of cycles: %d\n", cycles)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5),
				uint8(i%3))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(w, &anim) // NOTE: ignoring encoding errors
}

func doNothing(w http.ResponseWriter, r *http.Request) {}
