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
	Max           int
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var pgm PGM

	if !scanner.Scan() {
		return nil, fmt.Errorf("peu pas lire le magicnumber")
	}
	pgm.MagicNumber = scanner.Text()

	if pgm.MagicNumber != "P1" && pgm.MagicNumber != "P4" {
		return nil, fmt.Errorf("pas le bon format: %s", pgm.MagicNumber)
	}

	if !scanner.Scan() {
		return nil, fmt.Errorf("peu pas lire les dimention")
	}

	for scanner.Scan() {
		a := scanner.Text()
		if len(a) > 0 && a[0] == '#' {
			continue
		}
		fmt.Sscanf(a, "%d %d", &pgm.Width, &pgm.Height)
		break

	}
	pgm.Max, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid max value")
	}
	if pgm.MagicNumber == "P4" {
		pgm.Width *= 8
	}

	for scanner.Scan() {
		var binaryBits []uint8
		line := scanner.Text()
		tokens := strings.Fields(line)
		row := make([]uint8, pgm.Width)

		if pgm.MagicNumber == "P2" {
			for i, token := range tokens {
				if i >= pgm.Width {
					break
				}
				if token == "1" {
					row[i] = true
				} else if token == "0" {
					row[i] = false
				} else {
					return nil, fmt.Errorf("invalid character in data: %s", token)
				}
			}
		}
		if pgm.MagicNumber == "P2" {
			i := 0
			for _, token := range tokens {
				token = strings.TrimPrefix(token, "0x")
				for _, digit := range token {
					digitValue, err := strconv.ParseUint(string(digit), 16, 4)
					if err != nil {
						return nil, err
					}
					binaryDigits := strings.Split(fmt.Sprintf("%04b", digitValue), "")
					binaryBits = append(binaryBits, binaryDigits...)
				}

				if i >= pgm.Width {
					break
				}
				for _, value := range binaryBits {
					if value == "1" {
						row[i] = true
						i++
					} else if value == "0" {
						row[i] = false
						i++
					} else {
						return nil, fmt.Errorf("invalid character in data: %v", value)
					}
				}
			}
		}
		pgm.Data = append(pgm.Data, row)
	}
	return &pgm, nil
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
	fileName := "save.pgm"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", pgm.MagicNumber)
	fmt.Fprintf(file, "# saved file\n")
	fmt.Fprintf(file, "%d %d\n", pgm.Width, pgm.Height)
	
	for _, row := range pgm.Data {
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
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.Height; i++ {
		for j := 0; j < pgm.Width; j++ {
			pgm.Data[i][j] = ^pgm.Data[i][j]
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
    pgm.Max = int(maxValue)
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
    // Set a threshold value for conversion (adjust as needed)
    threshold := uint8(pgm.Max / 2)

    pbm := &PBM{
        Width:       pgm.Width,
        Height:      pgm.Height,
        MagicNumber: "P1",
    }
    for i := 0; i < pgm.Height; i++ {
        var row []bool
        for j := 0; j < pgm.Width; j++ {
            pixel := pgm.Data[i][j]
            if pixel > threshold {
                row = append(row, true)
            } else {
                row = append(row, false)
            }
        }
        pbm.Data = append(pbm.Data, row)
    }

    return pbm
}