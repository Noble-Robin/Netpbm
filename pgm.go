package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	Data          [][]uint8
	Width, Height int
	MagicNumber   string
	Max           uint
}
type PBM struct {
	Data          [][]bool
	Width, Height int
	MagicNumber   string
}


func ReadPGM(filename string) (*PGM, error) {
	var dimension string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var pgm PGM

	// Lecture de la première ligne pour obtenir le magic number
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P2" && line != "P5" {
		return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
	}
	pgm.MagicNumber = line

	// Lecture des dimensions
	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break

	}

	dimension = scanner.Text()
	res := strings.Split(dimension, " ")
	pgm.Height, _ = strconv.Atoi(res[0])
	pgm.Width, _ = strconv.Atoi(res[1])

	if !scanner.Scan() {
		return nil, fmt.Errorf("unable to read max value")
	}
	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse max value: %v", err)
	}
	pgm.Max = uint(maxValue)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		row := make([]uint8, pgm.Width)
		for i, token := range tokens {
			if i >= pgm.Width {
				break
			}
			value, err := strconv.ParseUint(token, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("caractère non valide dans les données : %s", token)
			}
			row[i] = uint8(value)
		}
		pgm.Data = append(pgm.Data, row)
	}
	return &PGM{
		Data:        pgm.Data,
		Width:       pgm.Width,
		Height:      pgm.Height,
		MagicNumber: pgm.MagicNumber,
		Max:         pgm.Max,
	}, nil
}
func (pgm *PGM) Size() (int, int) {
	return pgm.Width, pgm.Height
}
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.Data[x][y]
}
func (pgm *PGM) Set(x, y int, value uint8) {
	if x >= 0 && x < len(pgm.Data) && y >= 0 && y < len(pgm.Data[0]) {
		pgm.Data[x][y] = value
	}
}
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	header := fmt.Sprintf("%s\n%d %d\n%d\n", pgm.MagicNumber, pgm.Width, pgm.Height, pgm.Max)
	_, err = file.WriteString(header)
	if err != nil {
		return err
	}

	for _, row := range pgm.Data {
		for _, value := range row {
			_, err := fmt.Fprintf(file, "%d ", value)
			if err != nil {
				return err
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (pgm *PGM) Invert() {
	for i := 0; i < pgm.Height; i++ {
		for j := 0; j < pgm.Width; j++ {
			pgm.Data[i][j] = uint8(pgm.Max) - pgm.Data[i][j]
		}
	}
}
func (pgm *PGM) Flip() {
	for x := 0; x < pgm.Height; x++ {
		for i, j := 0, pgm.Width-1; i < j; i, j = i+1, j-1 {
			pgm.Data[x][i], pgm.Data[x][j] = pgm.Data[x][j], pgm.Data[x][i]
		}
	}
}
func (pgm *PGM) Flop() {
	for y := 0; y < pgm.Width; y++ {
		for i, j := 0, pgm.Height-1; i < j; i, j = i+1, j-1 {
			pgm.Data[i][y], pgm.Data[j][y] = pgm.Data[j][y], pgm.Data[i][y]
		}
	}
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.Max = uint(maxValue)
}
func (pgm *PGM) Rotate90CW() {
	rotatedData := make([][]uint8, pgm.Width)
	for i := range rotatedData {
		rotatedData[i] = make([]uint8, pgm.Height)
	}
	for i := 0; i < pgm.Height; i++ {
		for j := 0; j < pgm.Width; j++ {
			rotatedData[j][pgm.Height-1-i] = pgm.Data[i][j]
		}
	}
	pgm.Data = rotatedData
	pgm.Width, pgm.Height = pgm.Height, pgm.Width
}
func (pgm *PGM) ToPBM() *PBM {
    threshold := uint8(pgm.Max / 2) // Use midpoint as threshold
    pbmImage := &PBM{
        Width:       pgm.Width,
        Height:      pgm.Height,
        MagicNumber: "P1",
    }
    pbmImage.Data = make([][]bool, pbmImage.Height)
    for i := 0; i < pbmImage.Height; i++ {
        pbmImage.Data[i] = make([]bool, pbmImage.Width)
    }
    for i := 0; i < pgm.Height; i++ {
        for j := 0; j < pgm.Width; j++ {
            pbmImage.Data[i][j] = pgm.Data[i][j] > threshold
        }
    }

    return pbmImage
}