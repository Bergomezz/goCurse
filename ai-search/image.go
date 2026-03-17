package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/StephaneBunel/bresenham"
	"github.com/kmicki/apng"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const cellSize = 60

var (
	green = color.RGBA{G: 255, A: 255}
	darkGreen = color.RGBA{R: 1, G: 100,B: 32, A: 255}
	red =color.RGBA{R: 255, A: 255}
	yellow =color.RGBA{R: 255,G: 255,B: 101, A: 255}
	gray = color.RGBA{R: 125, G: 125, B: 125, A: 255}
	orange = color.RGBA{R: 255, G: 140, B: 25, A: 255}
	blue = color.RGBA{R: 14, G: 118, B: 173, A: 255}
)

func (g *Maze) OutputImage (fileName ...string) {
	fmt.Printf("Generate image %s...\n", fileName)
	width := cellSize * (g.Width -1)
	height := cellSize * g.Height

	var outFile = "image.png"
	if len(fileName) > 0 {
		outFile = fileName[0]
	}

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)

	for i, row := range g.Wall {
		for j, col := range row {
			p := Point{
				Row: i,
				Col: j,
			}
			if col.wall{
				g.drawSquare(col, p, img, color.Black, cellSize, j*cellSize, i*cellSize)
			}else if g.inSolution(p) {
				g.drawSquare(col, p, img, green, cellSize, j*cellSize, i*cellSize)
			}else if col.State.Row == g.Start.Row && col.State.Col == g.Start.Col {
				g.drawSquare(col, p, img, darkGreen, cellSize, j*cellSize, i*cellSize)
			}else if col.State.Row == g.Goal.Row && col.State.Col == g.Goal.Col {
				g.drawSquare(col, p, img, red, cellSize, j*cellSize, i*cellSize)
			}else if col.State == g.CurrentNode.State {
				g.drawSquare(col, p, img, orange, cellSize, j*cellSize, i*cellSize)
			}else if inExplored(Point{i, j}, g.Explored) {
				g.drawSquare(col, p, img, yellow, cellSize, j*cellSize, i*cellSize)
			}else {
				g.drawSquare(col, p, img, color.White, cellSize, j*cellSize, i*cellSize)
			}
		}
	}

	for i, _ := range g.Wall {
		bresenham.DrawLine(img, 0, i*cellSize, g.Width*cellSize, i*cellSize, gray)
	}

	for i := 0; i <= g.Width; i++ {
		bresenham.DrawLine(img, i*cellSize, 0, i*cellSize, g.Height*cellSize, gray)
	}

	f, _ := os.Create(outFile)

	_ = png.Encode(f, img)
}


func (g *Maze) drawSquare(col Wall, p Point, img *image.RGBA, c color.Color, size, x, y int) {
	patch := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(
		patch,
		patch.Bounds(),
		&image.Uniform{
			C: c,
		},
		image.Point{},
		draw.Src)
	if !col.wall {
		g.printLocation(p, color.Black, patch)
	}

	draw.Draw(img, image.Rect(x, y, x+size, y+size),patch, image.Point{}, draw.Src)
}

func (g *Maze) printLocation(p Point, c color.Color, patch *image.RGBA) {
	point := fixed.Point26_6{X: fixed.I(6), Y: fixed.I(40)}
	d := &font.Drawer{
		Dst: patch,
		Src: image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot: point,
	}

	d.DrawString(fmt.Sprintf("[%d %d]", p.Row, p.Col))
}

func (g *Maze) OutputAnimatedImage() {
	output := "./animation.png"
	files, _ := os.ReadDir("./tmp")
	var images []string
	var delays []int

	for _, file := range files {
		images = append(images, fmt.Sprintf("./tmp/%s", file.Name()))
		delays = append(delays, 30)
	}

	images = append(images, "./image.png")

	a := apng.APNG {
		Frames: make([]apng.Frame, len(images)),
	}

	out, _ := os.Create(output)
	defer out.Close()

	for i, s := range images {
		in, err := os.Open(s)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		defer in.Close()

		m, err := png.Decode(in)
		if err != nil {
			continue
		}
		a.Frames[i].Image = m
	}
	err := apng.Encode(out, a)
	if err != nil {
		log.Println(err)
	}
}