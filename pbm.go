package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	
	
)
type PBM struct{
    Data [][]bool
    Width, Height int
    MagicNumber string
}


func ReadPBM(filename string) (*PBM, error){
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	if !scanner.Scan() {
		return nil, fmt.Errorf("peu pas lire le magicnumber")
	}
	magicNumber := scanner.Text()

	if !scanner.Scan() {
		return nil, fmt.Errorf("peu pas lire les dimention")
	}
	var width, height int
	_, err = fmt.Sscanf(scanner.Text(), "%d %d", &width, &height)
	if err != nil {
		return nil, err
	}
	
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}

	for y := 0; y < height; y++ {
		if magicNumber == "P4" {
			if !scanner.Scan() {
				return nil, fmt.Errorf("peu pas lire la data de l'image")
			}
			row := scanner.Text()
			for x := 0; x < width; x++ {
				value, err := strconv.ParseInt(string(row[x/8]), 16, 64)
				if err != nil {
					return nil, err
				}
				data[y][x] = value&(1<<(7-x%8)) != 0
			}
		} else {
			for x := 0; x < width; x++ {
				var value int
				_, err := fmt.Fscanf(file, "%d", &value)
				if err != nil {
					return nil, err
				}
				data[y][x] = value == 1
			}
		}
	}

	return &PBM{Data: data, Width: width, Height: height, MagicNumber: magicNumber}, nil
}
