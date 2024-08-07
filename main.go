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
	GameHeight    = 240
	GameWidth     = 320
	DotSize       = 10
	ShowGridLines = false
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
		grid: NewGameGrid(GameWidth/(DotSize+1), GameHeight/(DotSize+1)),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	screenBuffer []byte
	grid         *GameGrid
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if x, y, ok := g.cursorGridPos(); ok {
			g.grid.Toggle(x, y)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.grid.Tick()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.buildGrid()
	g.buildLiveCells()
	g.buildCursorPos()

	screen.WritePixels(g.screenBuffer)
}

func (g *Game) buildGrid() {
	g.screenBuffer = make([]byte, GameWidth*GameHeight*4)

	if !ShowGridLines {
		return
	}

	for col := DotSize; col < GameWidth; col += DotSize + 1 {
		for row := 0; row < GameHeight; row++ {
			writePixel(g.screenBuffer, col, row, ColorWhite)
		}
	}

	for row := DotSize; row < GameHeight; row += DotSize + 1 {
		for col := 0; col < GameWidth; col++ {
			writePixel(g.screenBuffer, col, row, ColorWhite)
		}
	}
}

func (g *Game) buildLiveCells() {
	for i := 0; i < len(g.grid.Grid); i++ {
		if !g.grid.Grid[i] {
			continue
		}

		x, y := g.grid.Coords(i)
		x, y = screenCoords(x, y)

		for w := 0; w < DotSize; w++ {
			for h := 0; h < DotSize; h++ {
				writePixel(g.screenBuffer, x+w, h+y, ColorGreen)
			}
		}
	}
}

func (g *Game) buildCursorPos() {
	x, y, ok := g.cursorGridPos()
	if !ok {
		return
	}

	x, y = screenCoords(x, y)

	for w := 0; w < DotSize; w++ {
		for h := 0; h < DotSize; h++ {
			writePixel(g.screenBuffer, x+w, h+y, ColorRed)
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

func (g *Game) cursorGridPos() (x, y int, ok bool) {
	x, y = ebiten.CursorPosition()
	if x < 0 || y < 0 {
		return x, y, false
	}

	// xdiff := x % (DotSize + 1)
	// ydiff := y % (DotSize + 1)

	// if xdiff == 0 || ydiff == 0 {
	// 	return x, y, false
	// }

	x = int(math.Floor(float64(x) / (DotSize + 1)))
	y = int(math.Floor(float64(y) / (DotSize + 1)))

	if x >= g.grid.Width || y >= g.grid.Height {
		return x, y, false
	}

	return x, y, true
}

func screenCoords(x, y int) (int, int) {
	return x * (DotSize + 1), y * (DotSize + 1)
}
