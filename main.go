package main

import (
	"fmt"
	"os"
	path "path/filepath"

	"gradient-magic-wand/utils"
)

func prog(argc int, argv []string) int {

	if argc < 2 {
		fmt.Fprint(os.Stderr, "no parameter specified (use -h for help)\n")
		return 1
	}

	var filepath string = argv[1]

	img, err := utils.OpenImage(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open file \"%s\"\n%s\n", filepath, err.Error())
		return 1
	}

	saveFilepath := filepath[:len(filepath)-len(path.Ext(filepath))] + "_clean" + path.Ext(filepath)
	err = utils.SaveImage(saveFilepath, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not save file \"%s\"\n%s\n", saveFilepath, err.Error())
		return 1
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
		}
	}

	return 0
}

func main() {
	var argv []string = os.Args
	var argc = len(argv)

	os.Exit(prog(argc, argv))
}
