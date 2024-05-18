package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getParentPathFromFilename(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func loadImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var img image.Image = nil

	if strings.HasSuffix(filePath, ".png") {
		img, err = png.Decode(f)
	} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
		img, err = jpeg.Decode(f)
	} else {
		return img, errors.New("incorrect image format")
	}

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

func rgb_to_value(r uint8, g uint8, b uint8) uint8 {
	var max uint8 = 0

	if r > max {
		max = r
	}
	if g > max {
		max = g
	}
	if b > max {
		max = b
	}
	return max
}

func remove_bg(img image.NRGBA, min uint8, max uint8) image.NRGBA {
	size := img.Bounds().Size()
	var val uint64 = 0

	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			val = uint64(255 - rgb_to_value(uint8(r), uint8(g), uint8(b)))

			if val <= uint64(min) {
				val = 0
			} else {
				val = val*255/uint64(max-min) + uint64(min)
			}
			img.SetNRGBA(i, j, color.NRGBA{
				R: 0,
				G: 0,
				B: 0,
				A: uint8(val),
			})
		}
	}
	return img
}

func parseArguments() (string, uint8, uint8, error) {
	argv := os.Args
	argc := len(argv)

	if argc < 2 {
		return "", 0, 0, errors.New("too few arguments")
	}

	filename := argv[1]
	if !(strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") || strings.HasSuffix(filename, ".png")) {
		return "", 0, 0, errors.New("file not in a valid format")
	}

	var err error = nil
	var max int = 255
	var min int = 0
	for i := 2; i < argc; i += 2 {
		if argv[i][0] == '-' {
			if i == argc-1 {
				return "", 0, 0, errors.New("argument with no value")
			}
			if argv[i] == "-w" {
				max, err = strconv.Atoi(argv[i+1])
				if err != nil {
					return "", 0, 0, errors.New("white level not a number")
				}
			} else if argv[i] == "-b" {
				min, err = strconv.Atoi(argv[i+1])
				if err != nil {
					return "", 0, 0, errors.New("black level not a number")
				}
			} else {

			}
		} else {
			fmt.Printf("%s\n", argv[i])
			return "", 0, 0, errors.New("incorrect format")
		}
	}

	if max > 255 || max < 0 {
		return "", 0, 0, errors.New("white level not in range")
	}
	if min > 255 || min < 0 {
		return "", 0, 0, errors.New("black level not in range")
	}
	if min >= max {
		return "", 0, 0, errors.New("black level > white level")
	}

	return filename, uint8(min), uint8(max), nil
}

func main() {
	filename, min, max, err := parseArguments()
	if err != nil {
		log.Fatal(err)
	}

	img, err := loadImage(filename)
	if err != nil {
		log.Fatal(err)
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
