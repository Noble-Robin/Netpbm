package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	Data          [][]bool
	Width, Height int
	MagicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	var dimension string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var pbm PBM

	// Lecture de la première ligne pour obtenir le magic number
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P1" && line != "P4" {
		return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
	}
	pbm.MagicNumber = line

	// Lecture des dimensions
	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break

	}

	dimension = scanner.Text()
	res := strings.Split(dimension, " ")
	pbm.Height, _ = strconv.Atoi(res[0])
	pbm.Width, _ = strconv.Atoi(res[1])

	// Lecture des données binaires
	if pbm.MagicNumber == "P1" {
		pbm.Data = make([][]bool, pbm.Height)
		for i := range pbm.Data {
			pbm.Data[i] = make([]bool, pbm.Width)
		}
		for i := 0; i < pbm.Height; i++ {
			scanner.Scan()
			line := scanner.Text()
			hori := strings.Fields(line)
			for j := 0; j < pbm.Width; j++ {
				verti, _ := strconv.Atoi(hori[j])
				if verti == 1 {
					pbm.Data[i][j] = true
				}
			}
		}

	}
	// if pbm.MagicNumber =="P4"
	fmt.Printf("%+v\n", PBM{pbm.Data, pbm.Width, pbm.Height, pbm.MagicNumber})
	return &pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.Width, pbm.Height
}
func (pbm *PBM) At(x, y int) bool {
	return pbm.Data[x][y]
}
func (pbm *PBM) Set(x, y int, value bool) {
	if x >= 0 && x < len(pbm.Data) && y >= 0 && y < len(pbm.Data[0]) {
		pbm.Data[x][y] = value
	}
}
func (pbm *PBM) Save(filename string) error {
	fileName := "save.pbm"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", pbm.MagicNumber)
	fmt.Fprintf(file, "# saved file\n")
	fmt.Fprintf(file, "%d %d\n", pbm.Width, pbm.Height)
	for _, row := range pbm.Data {
		for _, pixel := range row {
			if pixel {
				fmt.Fprint(file, "1 ")

			} else {
				fmt.Fprint(file, "0 ")
			}
		}
		fmt.Fprintln(file)
	}

	fmt.Printf("File created: %s\n", fileName)
	return nil
}
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.Height; i++ {
		for j := 0; j < pbm.Width; j++ {
			pbm.Data[i][j] = !pbm.Data[i][j]
		}
	}
}
func (pbm *PBM) Flip() {
	for x := 0; x < pbm.Height; x++ {
		for i, j := 0, pbm.Width-1; i < j; i, j = i+1, j-1 {
			pbm.Data[x][i], pbm.Data[x][j] = pbm.Data[x][j], pbm.Data[x][i]
		}
	}
}
func (pbm *PBM) Flop() {
	for y := 0; y < pbm.Width; y++ {
		for i, j := 0, pbm.Height-1; i < j; i, j = i+1, j-1 {
			pbm.Data[i][y], pbm.Data[j][y] = pbm.Data[j][y], pbm.Data[i][y]
		}
	}
}
