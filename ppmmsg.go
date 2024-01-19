package netpbm

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sort"
)

// PPM is a structure to represent PPM images.
type PPM struct {
	Data        [][]Pixel
	Width, Height int
	MagicNumber string
	Max         uint
}

// Pixel represents a pixel with red (R), green (G), and blue (B) channels.
type Pixel struct {
	R, G, B uint8
}

// Point represents a point in the image.
type Point struct {
	X, Y int
}

// ReadPPM lit une image PPM depuis un fichier et renvoie une structure qui représente l'image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// lit le numero magique
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read magic number")
	}
	magicNumber := scanner.Text()

	if magicNumber != "P3" && magicNumber != "P6" {
		return nil, fmt.Errorf("unsupported PPM format: %s", magicNumber)
	}

	// passe les commentaires et les lignes vides
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' {
			break
		}
	}

	// lit la largeur et la hauteur
	if scanner.Err() != nil {
		return nil, fmt.Errorf("error reading dimensions line: %v", scanner.Err())
	}
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, fmt.Errorf("invalid dimensions line")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse width: %v", err)
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse height: %v", err)
	}

	// lit la valeur max
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read max value")
	}
	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse max value: %v", err)
	}

	// lit les données
	var data [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			var pixel Pixel
			if magicNumber == "P3" {
				// format ASCII
				if _, err := fmt.Fscanf(Scanner, "%d %d %d", &pixel.R, &pixel.G, &pixel.B); err != nil {
					return nil, fmt.Errorf("failed to parse pixel data: %v", err)
				}
			} else {
				// Format binaire (P6)
				var buf [3]byte
				if _, err := file.Read(buf[:]); err != nil {
					return nil, fmt.Errorf("failed to read pixel data: %v", err)
				}
				pixel.R, pixel.G, pixel.B = buf[0], buf[1], buf[2]
			}
			row = append(row, pixel)
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return &PPM{
		Data:        data,
		Width:       width,
		Height:      height,
		MagicNumber: magicNumber,
		Max:         uint(maxValue),
	}, nil
}