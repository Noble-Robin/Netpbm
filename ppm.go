package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PPM struct{
    Data [][]Pixel
    Width, Height int
    MagicNumber string
    Max int
}

type Pixel struct{
    R, G, B uint8
}

type Point struct{
    X, Y int
}

func (ppm *PPM) Size() (int, int) {
	return ppm.Width, ppm.Height
}
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.Data[x][y]
}
func (ppm *PPM) Set(x, y int, value Pixel) {
	if x >= 0 && x < len(ppm.Data) && y >= 0 && y < len(ppm.Data[0]) {
		ppm.Data[x][y] = value
	}
}
func (ppm *PPM) Save(filename string) error {
	fileName := "save.ppm"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", ppm.MagicNumber)
	fmt.Fprintf(file, "# saved file\n")
	fmt.Fprintf(file, "%d %d\n", ppm.Width, ppm.Height)
	
	for _, row := range ppm.Data {
		for _, pixel := range row {
			if pixel {
				fmt.Fprint(file, "1")
			} else {
				fmt.Fprint(file, "0")
			}
		}
		fmt.Fprintln(file)
	}

	fmt.Printf("File created: %s\n", fileName)
	return nil
}
func (ppm *PPM) Invert() {
    for i := 0; i < ppm.Height; i++ {
        for j := 0; j < ppm.Width; j++ {
            // Invert the RGB components of the current pixel
            ppm.Data[i][j].R = 255 - ppm.Data[i][j].R
            ppm.Data[i][j].G = 255 - ppm.Data[i][j].G
            ppm.Data[i][j].B = 255 - ppm.Data[i][j].B
        }
    }
}
func (ppm *PPM) Flip() {
	for x := 0; x < ppm.Height; x++ {
		for i, j := 0, ppm.Width-1; i < j; i, j = i+1, j-1 {
			ppm.Data[x][i], ppm.Data[x][j] = ppm.Data[x][j], ppm.Data[x][i]
		}
	}
}
func (ppm *PPM) Flop() {
	for y := 0; y < ppm.Width; y++ {
		for i, j := 0, ppm.Height-1; i < j; i, j = i+1, j-1 {
			ppm.Data[i][y], ppm.Data[j][y] = ppm.Data[j][y], ppm.Data[i][y]
		}
	}
}
func (ppm *PPM) SetMaxValue(maxValue uint8) {
    ppm.Max = int(maxValue)
}
func (ppm *PPM) Rotate90CW() {
    rotatedData := make([][]Pixel, ppm.Width)
    for i := range rotatedData {
        rotatedData[i] = make([]Pixel, ppm.Height)
    }
    for i := 0; i < ppm.Height; i++ {
        for j := 0; j < ppm.Width; j++ {
            rotatedData[j][ppm.Height-1-i] = ppm.Data[i][j]
        }
    }
    ppm.Data = rotatedData
    ppm.Width, ppm.Height = ppm.Height, ppm.Width
}
func (ppm *PPM) ToPGM() *PGM {
    pgm := &PGM{
        Width:       ppm.Width,
        Height:      ppm.Height,
        MagicNumber: "P2",
        Max:         255,
    }
    for i := 0; i < ppm.Height; i++ {
        var row []uint8
        for j := 0; j < ppm.Width; j++ {
            average := uint8((uint16(ppm.Data[i][j].R) + uint16(ppm.Data[i][j].G) + uint16(ppm.Data[i][j].B)) / 3)
            row = append(row, average)
        }
        pgm.Data = append(pgm.Data, row)
    }

    return pgm
}
func (ppm *PPM) ToPBM() *PBM {
    pbm := &PBM{
        Width:       ppm.Width,
        Height:      ppm.Height,
        MagicNumber: "P1",
    }

    threshold := uint8((ppm.Max / 2) + 1)

    for i := 0; i < ppm.Height; i++ {
        var row []bool
        for j := 0; j < ppm.Width; j++ {
           
            average := (uint16(ppm.Data[i][j].R) + uint16(ppm.Data[i][j].G) + uint16(ppm.Data[i][j].B)) / 3
            if average > uint16(threshold) {
                row = append(row, true)
            } else {
                row = append(row, false)
            }
        }
        pbm.Data = append(pbm.Data, row)
    }

    return pbm
}
