// Copyright (C) 2015 The Protogalaxy Project
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package tictactoe

const GridSize int = 3

type gameGrid struct {
	grid []Mark
}

func newGameGrid() *gameGrid {
	return &gameGrid{
		grid: make([]Mark, GridSize*GridSize),
	}
}

func (g *gameGrid) set(x, y int, m Mark) {
	g.grid[y*GridSize+x] = m
}

func (g *gameGrid) isEmpty(x, y int) bool {
	return g.get(x, y) == Mark_EMPTY
}

func (g *gameGrid) get(x, y int) Mark {
	return g.grid[y*GridSize+x]
}

func (g *gameGrid) checkMarksHorizontal(y int) bool {
	first := g.get(0, y)
	return first == g.get(1, y) && first == g.get(2, y)
}

func (g *gameGrid) checkDiagonalDown() bool {
	first := g.get(0, 0)
	return first == g.get(1, 1) && first == g.get(2, 2) && first != Mark_EMPTY
}

func (g *gameGrid) checkDiagonalUp() bool {
	first := g.get(0, 2)
	return first == g.get(1, 1) && first == g.get(2, 0) && first != Mark_EMPTY
}

func (g *gameGrid) checkMarksVertical(x int) bool {
	first := g.get(x, 0)
	return first == g.get(x, 1) && first == g.get(x, 2)
}

func (g *gameGrid) isFull() bool {
	for _, m := range g.grid {
		if m == Mark_EMPTY {
			return false
		}
	}
	return true
}

func (g *gameGrid) clone() *gameGrid {
	grid := make([]Mark, GridSize*GridSize)
	copy(grid, g.grid)
	return &gameGrid{grid: grid}
}

func (g *gameGrid) coordinatesValid(x, y int) bool {
	return validateIndex(x) && validateIndex(y)
}

func validateIndex(x int) bool {
	return x >= 0 && x < GridSize
}
