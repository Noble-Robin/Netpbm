package main

import (
    "fmt"
    "netpbm/projet"
)

func main() {
	img, err := ReadPBM("image.pbm")
	if err != nil {
		fmt.Println("Error reading PBM file:", err)
		return
	}

	fmt.Println("Width:", img.width)
	fmt.Println("Height:", img.height)
	fmt.Println("Magic Number:", img.magicNum)
	fmt.Println("Data:", img.data)
}