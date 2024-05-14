package main

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func getParentPathFromFilename(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func loadPNG(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	return img, err
}

func savePNG(img image.NRGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, &img)
	if err != nil {
		return err
	}

	return nil
}

func clamp(min int, x int, max int) int {
	if x > max {
		return max
	}
	if x < min {
		return min
	}
	return x
}

func remove_bg(img image.NRGBA, min int, max int) image.NRGBA {
	size := img.Bounds().Size()

	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			ur := uint8(r)
			ug := uint8(g)
			ub := uint8(b)
			avg := (int(ur) + int(ug) + int(ub)) / 3

			img.SetNRGBA(i, j, color.NRGBA{
				R: ur,
				G: ug,
				B: ub,
				A: uint8(clamp(0, max-avg, 255)),
			})
		}
	}

	return img
}

func main() {
	argv := os.Args
	argc := len(os.Args)

	if argc != 4 {
		log.Fatal(errors.New("missing argument(s)"))
	}

	filename := argv[1]
	img, err := loadPNG(filename)
	if err != nil {
		log.Fatal(err)
	}
	min, err := strconv.Atoi(argv[2])
	if err != nil {
		log.Fatal(errors.New("min not a number"))
	}
	max, err := strconv.Atoi(argv[3])
	if err != nil {
		log.Fatal(errors.New("max not a number"))
	}

	bounds := img.Bounds()
	imgRGBA := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), img, bounds.Min, draw.Src)

	outputImg := remove_bg(*imgRGBA, min, max)

	savePNG(outputImg, getParentPathFromFilename(filename)+".clean.png")

	if err != nil {
		log.Fatal(err)
	}
}
