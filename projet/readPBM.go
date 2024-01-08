package projet

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"Netpbm/struct"
	
)

func ReadPBM(filename string) (*PBM, error){
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the header
	if !scanner.Scan() {
		return nil, fmt.Errorf("could not read magic number")
	}
	magicNum := scanner.Text()

	if !scanner.Scan() {
		return nil, fmt.Errorf("could not read width and height")
	}
	var width, height int
	_, err = fmt.Sscanf(scanner.Text(), "%d %d", &width, &height)
	if err != nil {
		return nil, err
	}

	// Initialize the image data
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}

	// Read the image data
	for y := 0; y < height; y++ {
		if magicNum == "P4" {
			if !scanner.Scan() {
				return nil, fmt.Errorf("could not read image data")
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
				_, err := fmt.Fscanf(scanner, "%d", &value)
				if err != nil {
					return nil, err
				}
				data[y][x] = value == 1
			}
		}
	}

	return &PBM{data: data, width: width, height: height,}, nil
}