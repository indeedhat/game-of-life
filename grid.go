package main

import (
	"log"
)

type GameGrid struct {
	Width  int
	Height int
	Grid   []bool
}

func NewGameGrid(w, h int) *GameGrid {
	return &GameGrid{
		Width:  w,
		Height: h,
		Grid:   make([]bool, w*h),
	}
}

func (g *GameGrid) Index(x, y int) int {
	return y*g.Width + x
}

func (g *GameGrid) Coords(i int) (x, y int) {
	x = i % g.Width
	y = (i - x) / g.Width
	return x, y
}

func (g *GameGrid) Toggle(x, y int) {
	i := g.Index(x, y)
	g.Grid[i] = !g.Grid[i]

	log.Print(x, y, i)
}

func (g *GameGrid) Tick() {
	log.Print("tick ", g.Width, g.Height)
	frame := make([]bool, len(g.Grid))

	for i := range g.Grid {
		switch g.neighbourCount(i) {
		case 2:
			frame[i] = g.Grid[i]
		case 3:
			frame[i] = true
		}
	}

	g.Grid = frame
}

func (g *GameGrid) neighbourCount(i int) int {
	var count int
	if g.Grid[i] {
		log.Print("filled: ", i)
	}

	for y := -g.Width; y <= g.Width; y += g.Width {
		if i+y < 0 || i+y >= len(g.Grid) {
			continue
		}

		for x := -1; x <= 1; x++ {
			if x == 0 && y == 0 {
				continue
			}

			// NB: there is no need to check the y positions here as they will just be checked after the increments
			if x == -1 && i%g.Width == 0 ||
				x == 1 && i%g.Width == g.Width-1 {
				continue
			}

			neighbour := i + x + y

			if neighbour < 0 || neighbour >= len(g.Grid) {
				continue
			}

			if g.Grid[neighbour] {
				count++
			}
		}
	}

	return count
}
