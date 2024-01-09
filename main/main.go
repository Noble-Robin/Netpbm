package main

import (
    "fmt"
    "github.com/Noble-Robin/Netpbm"
	)

	
func main() {
	filename := "testP1.pbm"
	pbm, err := netpbm.ReadPBM(filename)
	if err != nil {
		fmt.Println("impossible de lire le fichier", err)
		return
	}

	fmt.Println("Width:", pbm.Width)
	fmt.Println("Height:", pbm.Height)
	fmt.Println("Magic Number:", pbm.MagicNumber)
	fmt.Println("Data:", pbm.Data)

}
