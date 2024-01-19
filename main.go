func main() {
	filename := "testP2.pgm"
	pgm, err := ReadPGM(filename)
	if err != nil {
		fmt.Println("impossible de lire le fichier", err)
		return
	}
	fmt.Println(pgm.Data)
	pgm.Save("savep2.pgm")
	// for x := 0; x < pgm.Height; x++ {
	// 	fmt.Println()
	// 	for y := 0; y < pgm.Width; y++ {
	// 		if pgm.Data[x][y] == 0 {
	// 			fmt.Print("□")
	// 		} else {
	// 			fmt.Print("■")
	// 		}

	// 	}
	// }
}

func main() {
	filename := "testP1.pbm"
	pbm, err := ReadPBM(filename)
	if err != nil {
		fmt.Println("impossible de lire le fichier", err)
		return
	}
	// pbm.Set(pbm.Width,pbm.Height,true)

	// pbm.Flop()
	// pbm.Flip()
	// pbm.Invert()
	for x := 0; x < pbm.Height; x++ {
		 fmt.Println()
		for y := 0; y < pbm.Width; y++ {
			if pbm.Data[x][y] == true {
				fmt.Print("□")
			} else {
				fmt.Print("■")
			}

		}
	}
	fmt.Println()
	pbm.Save(filename)
}

func main() {
	filename := "testP3.ppm"
	// ppm, err := ReadPPM(filename)
	// if err != nil {
	// 	fmt.Println("impossible de lire le fichier", err)
	// 	return

	// }
	fmt.Println(ReadPPM(filename))
	// ppm.Save("savep2.pgm")
	// for x := 0; x < ppm.Height; x++ {
	// 	fmt.Println()
	// 	for y := 0; y < ppm.Width; y++ {
	// 		if ppm.Data[x][y] == 0 {
	// 			fmt.Print("□")
	// 		} else {
	// 			fmt.Print("■")
	// 		}
	// 	}
	// }
}