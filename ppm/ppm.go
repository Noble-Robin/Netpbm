package ppm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type PPM struct {
	Data          [][]Pixel
	Width, Height int
	MagicNumber   string
	Max           uint
}
type Pixel struct {
	R, G, B uint8
}
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
type Point struct {
	X, Y int
}
func ReadPPM(filename string) (*PPM, error) {
	var dimension string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ppm PPM
	var pixel Pixel

	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P3" && line != "P6" {
		return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
	}
	ppm.MagicNumber = line

	// Lecture des dimensions
	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break
	}
	dimension = scanner.Text()
	res := strings.Split(dimension, " ")
	ppm.Height, _ = strconv.Atoi(res[0])
	ppm.Width, _ = strconv.Atoi(res[1])

	if !scanner.Scan() {
		return nil, fmt.Errorf("unable to read max value")
	}
	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse max value: %v", err)
	}
	ppm.Max = uint(maxValue)
	if scanner.Scan() {
		fmt.Fscanf(strings.NewReader(scanner.Text()), "%d %d %d", &pixel.R, &pixel.G, &pixel.B)
	}
	for y := 0; y < ppm.Height; y++ {
		var row []Pixel
		for x := 0; x < ppm.Width; x++ {
			if ppm.MagicNumber == "P3" {
				if _, err := fmt.Fscanf(strings.NewReader(scanner.Text()), "%d %d %d", &pixel.R, &pixel.G, &pixel.B); err != nil {
					return nil, fmt.Errorf("échec de l'analyse des données des pixels : %v", err)
				}
			}
			row = append(row, pixel)
		}
		ppm.Data = append(ppm.Data, row)
	}

	return &PPM{
		Data:        ppm.Data,
		Width:       ppm.Width,
		Height:      ppm.Height,
		MagicNumber: ppm.MagicNumber,
		Max:         ppm.Max,
	}, nil
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
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write magic number, width, height, and max value to the file
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.MagicNumber, ppm.Width, ppm.Height, ppm.Max)

	// Write pixel data to the file
	for y := 0; y < ppm.Height; y++ {
		for x := 0; x < ppm.Width; x++ {
			pixel := ppm.Data[y][x]
			fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
		}
		fmt.Fprintln(file) // Newline at the end of each row
	}

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
	ppm.Max = uint(maxValue)
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
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := int(math.Abs(float64(p2.X - p1.X)))
	dy := int(math.Abs(float64(p2.Y - p1.Y)))
	var sx, sy int

	if p1.X < p2.X {
		sx = 1
	} else {
		sx = -1
	}

	if p1.Y < p2.Y {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		ppm.data[p1.Y][p1.X] = color

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width - 1, p1.Y + height - 1}

	// Draw top and bottom edges
	for x := p1.X; x <= p2.X; x++ {
		ppm.Data[p1.Y][x] = color
		ppm.Data[p2.Y][x] = color
	}

	// Draw left and right edges (excluding corners to avoid duplicates)
	for y := p1.Y + 1; y < p2.Y; y++ {
		ppm.Data[y][p1.X] = color
		ppm.Data[y][p2.X] = color
	}
}
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
    p2 := Point{p1.X + width - 1, p1.Y + height - 1}

    for y := p1.Y; y <= p2.Y; y++ {
        for x := p1.X; x <= p2.X; x++ {
            ppm.Data[y][x] = color
        }
    }
}
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
    x := radius
    y := 0
    err := 0

    for x >= y {
        if x+center.X >= 0 && x+center.X < ppm.Width && y+center.Y >= 0 && y+center.Y < ppm.Height {
            ppm.Data[y+center.Y][x+center.X] = color
        }
        if y+center.X >= 0 && y+center.X < ppm.Width && x+center.Y >= 0 && x+center.Y < ppm.Height {
            ppm.Data[x+center.Y][y+center.X] = color
        }
        if -x+center.X >= 0 && -x+center.X < ppm.Width && y+center.Y >= 0 && y+center.Y < ppm.Height {
            ppm.Data[y+center.Y][-x+center.X] = color
        }
        if -y+center.X >= 0 && -y+center.X < ppm.Width && x+center.Y >= 0 && x+center.Y < ppm.Height {
            ppm.Data[x+center.Y][-y+center.X] = color
        }
        if -x+center.X >= 0 && -x+center.X < ppm.Width && -y+center.Y >= 0 && -y+center.Y < ppm.Height {
            ppm.Data[-y+center.Y][-x+center.X] = color
        }
        if -y+center.X >= 0 && -y+center.X < ppm.Width && -x+center.Y >= 0 && -x+center.Y < ppm.Height {
            ppm.Data[-x+center.Y][-y+center.X] = color
        }
        if x+center.X >= 0 && x+center.X < ppm.Width && -y+center.Y >= 0 && -y+center.Y < ppm.Height {
            ppm.Data[-y+center.Y][x+center.X] = color
        }
        if y+center.X >= 0 && y+center.X < ppm.Width && -x+center.Y >= 0 && -x+center.Y < ppm.Height {
            ppm.Data[-x+center.Y][y+center.X] = color
        }

        y++
        if err <= 0 {
            err += 2*y + 1
        }
        if err > 0 {
            x--
            err -= 2*x + 1
        }
    }
}
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
    x := radius
    y := 0
    err := 0

    for x >= y {
        for i := -x; i <= x; i++ {
            if center.Y+y >= 0 && center.Y+y < ppm.Height && center.X+i >= 0 && center.X+i < ppm.Width {
                ppm.Data[center.Y+y][center.X+i] = color
            }
            if center.Y-y >= 0 && center.Y-y < ppm.Height && center.X+i >= 0 && center.X+i < ppm.Width {
                ppm.Data[center.Y-y][center.X+i] = color
            }
        } 
        for i := -y; i <= y; i++ {
            if center.Y+i >= 0 && center.Y+i < ppm.Height && center.X+x >= 0 && center.X+x < ppm.Width {
                ppm.Data[center.Y+i][center.X+x] = color
            }
            if center.Y+i >= 0 && center.Y+i < ppm.Height && center.X-x >= 0 && center.X-x < ppm.Width {
                ppm.Data[center.Y+i][center.X-x] = color
            }
        }
        y++
        if err <= 0 {
            err += 2*y + 1
        }
        if err > 0 {
            x--
            err -= 2*x + 1
        }
    }
}
func drawHorizontalLine(ppm *PPM, x1, x2, y int, color Pixel) {
    if x1 > x2 {
        x1, x2 = x2, x1
    }

    for x := x1; x <= x2; x++ {
        if x >= 0 && x < ppm.Width && y >= 0 && y < ppm.Height {
            ppm.Data[y][x] = color
        }
    }
}
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
    if p1.Y > p2.Y {
        p1, p2 = p2, p1
    }
    if p2.Y > p3.Y {
        p2, p3 = p3, p2
    }
    if p1.Y > p2.Y {
        p1, p2 = p2, p1
    }
    slope1 := float64(p2.X-p1.X) / float64(p2.Y-p1.Y)
    slope2 := float64(p3.X-p1.X) / float64(p3.Y-p1.Y)
    slope3 := float64(p3.X-p2.X) / float64(p3.Y-p2.Y)
    x1 := float64(p1.X)
    x2 := float64(p1.X)
    x3 := float64(p2.X)
    for y := p1.Y; y <= p2.Y; y++ {
        drawHorizontalLine(ppm, int(x1), int(x2), y, color)
        x1 += slope1
        x2 += slope2
    }
    x2 = float64(p2.X)
    for y := p2.Y + 1; y <= p3.Y; y++ {
        drawHorizontalLine(ppm, int(x2), int(x3), y, color)
        x2 += slope3
        x3 += slope2
    }
}
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
    if p1.Y > p2.Y {
        p1, p2 = p2, p1
    }
    if p2.Y > p3.Y {
        p2, p3 = p3, p2
    }
    if p1.Y > p2.Y {
        p1, p2 = p2, p1
    }
    slope1 := float64(p2.X-p1.X) / float64(p2.Y-p1.Y)
    slope2 := float64(p3.X-p1.X) / float64(p3.Y-p1.Y)
    slope3 := float64(p3.X-p2.X) / float64(p3.Y-p2.Y)
    x1 := float64(p1.X)
    x2 := float64(p1.X)
    x3 := float64(p2.X)  
    for y := p1.Y; y <= p2.Y; y++ {
        drawHorizontalLine(ppm, int(x1), int(x2), y, color)
        x1 += slope1
        x2 += slope2
    }
    x2 = float64(p2.X)
    for y := p2.Y + 1; y <= p3.Y; y++ {
        drawHorizontalLine(ppm, int(x2), int(x3), y, color)
        x2 += slope3
        x3 += slope2
    }
}
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
    if len(points) < 3 {
        return
    }
    sort.Slice(points, func(i, j int) bool {
        return points[i].Y < points[j].Y
    })
    for y := points[0].Y; y <= points[len(points)-1].Y; y++ {
        intersections := []int{}
        for i := 0; i < len(points); i++ {
            j := (i + 1) % len(points)
            if (points[i].Y <= y && points[j].Y > y) || (points[j].Y <= y && points[i].Y > y) {
                x := int(float64(points[i].X) + float64(y-points[i].Y)/float64(points[j].Y-points[i].Y)*float64(points[j].X-points[i].X))
                intersections = append(intersections, x)
            }
        }
        sort.Ints(intersections)
        for i := 0; i < len(intersections); i += 2 {
            x1 := intersections[i]
            x2 := intersections[i+1]
            drawHorizontalLine(ppm, x1, x2, y, color)
        }
    }
}
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
    if len(points) < 3 {
        return
    }
    sort.Slice(points, func(i, j int) bool {
        return points[i].Y < points[j].Y
    })
    for y := points[0].Y; y <= points[len(points)-1].Y; y++ {
        intersections := []int{}
        for i := 0; i < len(points); i++ {
            j := (i + 1) % len(points)
            if (points[i].Y <= y && points[j].Y > y) || (points[j].Y <= y && points[i].Y > y) {
                x := int(float64(points[i].X) + float64(y-points[i].Y)/float64(points[j].Y-points[i].Y)*float64(points[j].X-points[i].X))
                intersections = append(intersections, x)
            }
        }
        sort.Ints(intersections)
        for i := 0; i < len(intersections); i += 2 {
            x1 := intersections[i]
            x2 := intersections[i+1]
            drawHorizontalLine(ppm, x1, x2, y, color)
        }
    }
}
