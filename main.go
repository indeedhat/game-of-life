package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ColorWhite = Pixel{0xff, 0xff, 0xff, 0xff}
	ColorRed   = Pixel{0xff, 0x0, 0x0, 0xff}
	ColorGreen = Pixel{0x0, 0xff, 0x0, 0xff}
)

const (
	GameHeight = 240
	GameWidth  = 320
	DotSize    = 9
)

type Pixel struct {
	R byte
	G byte
	B byte
	A byte
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	game := &Game{
		selectedGrid: make(Grid, GameWidth/5*GameHeight/5),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	grid         []byte
	selectedGrid Grid
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if x, y, ok := cursorGridPos(); ok {
			i := g.selectedGrid.coordsToIndex(x, y)
			log.Print(x, y, i)
			log.Print(g.selectedGrid.indexToCoords(i))
			g.selectedGrid[i] = !g.selectedGrid[i]

			log.Print(len(g.selectedGrid))
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.buildGrid()
	g.buildSelectedDots()
	g.buildCursorPos()

	screen.WritePixels(g.grid)
}

func (g *Game) buildGrid() {
	g.grid = make([]byte, GameWidth*GameHeight*4)

	for col := DotSize; col < GameWidth; col += DotSize + 1 {
		for row := 0; row < GameHeight; row++ {
			writePixel(g.grid, col, row, ColorWhite)
		}
	}

	for row := DotSize; row < GameHeight; row += DotSize + 1 {
		for col := 0; col < GameWidth; col++ {
			writePixel(g.grid, col, row, ColorWhite)
		}
	}
}

func (g *Game) buildSelectedDots() {
	for i := 0; i < len(g.selectedGrid); i++ {
		if !g.selectedGrid[i] {
			continue
		}

		x, y := g.selectedGrid.indexToCoords(i)
		x, y = g.selectedGrid.toScreen(x, y)

		for w := 0; w < DotSize; w++ {
			for h := 0; h < DotSize; h++ {
				writePixel(g.grid, x+w, h+y, ColorGreen)
			}
		}
	}
}

func (g *Game) buildCursorPos() {
	x, y, ok := cursorGridPos()
	if !ok {
		return
	}

	x, y = g.selectedGrid.toScreen(x, y)

	for w := 0; w < DotSize; w++ {
		for h := 0; h < DotSize; h++ {
			writePixel(g.grid, x+w, h+y, ColorRed)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func writePixel(grid []byte, x, y int, color Pixel) {
	i := x * 4
	i += y * GameWidth * 4

	grid[i] = color.R
	grid[i+1] = color.G
	grid[i+2] = color.B
	grid[i+3] = color.A
}

func cursorGridPos() (x, y int, ok bool) {
	x, y = ebiten.CursorPosition()
	if x < 0 || y < 0 {
		return x, y, false
	}

	xdiff := x % (DotSize + 1)
	ydiff := y % (DotSize + 1)

	if xdiff == 0 || ydiff == 0 {
		return x, y, false
	}

	x = int(math.Floor(float64(x) / (DotSize + 1)))
	y = int(math.Floor(float64(y) / (DotSize + 1)))

	return x, y, true
}

type Grid []bool

func (Grid) toScreen(x, y int) (int, int) {
	if x > 0 {
		x = x * (DotSize + 1)
	}
	if y > 0 {
		y = y * (DotSize + 1)
	}

	return x, y
}

func (Grid) coordsToIndex(x, y int) int {
	return x + y*GameWidth/(DotSize+1)
}

func (Grid) indexToCoords(i int) (x, y int) {
	x = i % (GameWidth / (DotSize + 1))
	y = int((i - x) / (GameWidth / (DotSize + 1)))

	return x, y
}
