package main

import (
	"flag"
	"fmt"
)

var board = [][]uint8{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func run() {
	var iterations uint

	flag.UintVar(&iterations, "n", 1, "The number of iterations to run")
	flag.Parse()

	for i := 0; i < int(iterations); i++ {
		iterate(&board)
	}

	for _, row := range board {
		fmt.Println(row)
	}
}

// iterate runs a single step of the rules
func iterate(state *[][]uint8) {
	newState := make([][]uint8, len(*state))

	for y, row := range *state {
		newState[y] = make([]uint8, len(row))

		for x, cell := range row {
			count := neighbourCount(x, y, *state)

			switch true {
			case count < 2:
				newState[y][x] = 0
			case count > 2:
				newState[y][x] = 1
			default:
				newState[y][x] = cell
			}
		}
	}

	*state = newState
}

// neighbourCount of the cell at the given coords
func neighbourCount(x, y int, state [][]uint8) int {
	var count int
	width := len(state)
	height := len(state[0])

	for i := -1; i <= 1; i++ {
		if y+i < 0 || y+i >= height {
			continue
		}

		for j := -1; j <= 1; j++ {
			if x+j < 0 || x+j >= width {
				continue
			}

			if i == 0 && j == 0 {
				continue
			}

			if count >= 3 {
				continue
			}

			count += int(state[y+i][x+j])
		}
	}

	return count
}
