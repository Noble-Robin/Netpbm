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
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var pbm PBM

	if !scanner.Scan() {
		return nil, fmt.Errorf("peu pas lire le magicnumber")
	}
	pbm.MagicNumber = scanner.Text()

	if !scanner.Scan() {
		return nil, fmt.Errorf("peu pas lire les dimention")
	}

	for scanner.Scan() {
		a := scanner.Text()
		if len(a) > 0 && a[0] == '#' {
			continue
		}
		fmt.Sscanf(a, "%d %d", &pbm.Width, &pbm.Height)
		break

	}

	for scanner.Scan() {
		var binaryBits []string
		line := scanner.Text()
		tokens := strings.Fields(line)
		row := make([]bool, pbm.Width)

		if pbm.MagicNumber == "P1" {
			for i, token := range tokens {
				if i >= pbm.Width {
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
		if pbm.MagicNumber == "P4" {
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

				if i >= pbm.Width {
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
		pbm.Data = append(pbm.Data, row)
	}
	return &pbm, nil
}

func main() {
	pbm, err := ReadPBM("test.pbm")
	if err != nil {
		fmt.Println("impossible de lire le fichier", err)
		return
	}
	fmt.Println("Width:", pbm.Width)
	fmt.Println("Height:", pbm.Height)
	fmt.Println("Magic Number:", pbm.MagicNumber)
	fmt.Println("Data:", pbm.Data)
}
