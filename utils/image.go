package utils

import (
	"errors"
	"image"
	"image/png"
	"os"
	"strings"
)

func OpenImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

func SaveImage(filePath string, img image.Image) error {
	if !strings.HasSuffix(filePath, ".png") {
		return errors.New("file musts be saved as a png")
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}
