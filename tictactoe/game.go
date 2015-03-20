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

import (
	"errors"
	"time"
)

type GameID string

type game struct {
	ID            GameID
	Grid          *gameGrid
	CurrentPlayer int
	PlayerList    []string
	Players       map[string]Mark
	Winner        *Winner
	TurnNumber    int
	TurnTimestamp int64
}

func newGame(ID GameID, playerOne, playerTwo string) *game {
	return &game{
		ID:         ID,
		Grid:       newGameGrid(),
		PlayerList: []string{playerOne, playerTwo},
		Players: map[string]Mark{
			playerOne: Mark_X,
			playerTwo: Mark_Y,
		},
	}
}

var (
	ErrInvalidMove     = errors.New("invalid move")
	ErrNotActivePlayer = errors.New("not active player")
	ErrInvalidMoveID   = errors.New("invalid move id")
)

func (g *game) activePlayer() string {
	if g.isFinished() {
		return ""
	}
	return g.PlayerList[g.CurrentPlayer]
}

func (g *game) updateActivePlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.PlayerList)
}

func (g *game) checkWinner(userID string, x, y int) {
	w := &Winner{}
	var noWinner bool
	if g.Grid.checkMarksHorizontal(y) {
		w.Locations = append(w.Locations, &Winner_Location{
			Direction: Winner_Location_HORIZONTAL,
			Position:  int32(y),
		})
		w.UserId = userID
	} else if g.Grid.checkMarksVertical(x) {
		w.Locations = append(w.Locations, &Winner_Location{
			Direction: Winner_Location_VERTICAL,
			Position:  int32(x),
		})
		w.UserId = userID
	} else if g.Grid.checkDiagonalUp() {
		w.Locations = append(w.Locations, &Winner_Location{
			Direction: Winner_Location_DIAGONAL_UP,
		})
		w.UserId = userID
	} else if g.Grid.checkDiagonalDown() {
		w.Locations = append(w.Locations, &Winner_Location{
			Direction: Winner_Location_DIAGONAL_DOWN,
		})
		w.UserId = userID
	} else if g.Grid.isFull() {
		w.Draw = true
	} else {
		noWinner = true
	}
	if !noWinner {
		g.Winner = w
	}
}

func (g *game) isFinished() bool {
	return g.Winner != nil
}

func (g *game) isDraw() bool {
	return g.Winner.Draw
}

func (g *game) winner() *Winner {
	return g.Winner
}

func (g *game) placeMark(userID string, moveID int64, x, y int) error {
	if g.activePlayer() != userID {
		return ErrNotActivePlayer
	} else if moveID != g.lastMoveID() {
		return ErrInvalidMoveID
	} else if !g.Grid.coordinatesValid(x, y) || !g.Grid.isEmpty(x, y) {
		return ErrInvalidMove
	}
	g.Grid.set(x, y, g.Players[userID])
	g.checkWinner(userID, x, y)
	g.updateActivePlayer()
	g.TurnNumber += 1
	g.TurnTimestamp = time.Now().UnixNano()
	return nil
}

func (g *game) lastMoveID() int64 {
	return (g.TurnTimestamp << 16) & int64(g.TurnNumber)
}

func (g *game) validMoves() []*MoveRange {
	occupied := g.Grid.clone()
	validMoves := make([]*MoveRange, 0)
	for x := 0; x < GridSize; x++ {
		for y := 0; y < GridSize; y++ {
			if !occupied.isEmpty(x, y) {
				continue
			}
			moveRange := findValidMoveRange(occupied, x, y)
			validMoves = append(validMoves, moveRange)
			markOccupied(occupied, moveRange)
		}
	}
	return validMoves
}

func markOccupied(g *gameGrid, m *MoveRange) {
	if m.ToX == 0 && m.ToY == 0 {
		g.set(int(m.FromX), int(m.FromY), Mark_X)
		return
	}
	for x := m.FromX; x <= m.ToX; x++ {
		for y := m.FromY; y <= m.ToY; y++ {
			g.set(int(x), int(y), Mark_X)
		}
	}
}

func findValidMoveRange(g *gameGrid, x, y int) *MoveRange {
	endY := searchValidVertical(g, x, y)
	endX := extendHorizontally(g, x, y, endY)
	m := &MoveRange{
		FromX: int32(x),
		FromY: int32(y),
		ToY:   int32(endY),
		ToX:   int32(endX),
	}
	if m.FromX == m.ToX && m.FromY == m.ToY {
		m.ToX = 0
		m.ToY = 0
	}
	return m
}

func searchValidVertical(g *gameGrid, posX, posY int) int {
	y := posY
	for ; y < GridSize; y++ {
		if !g.isEmpty(posX, y) {
			break
		}
	}
	y -= 1
	return y
}

func extendHorizontally(g *gameGrid, posX, posY, endY int) int {
	x := posX
	for ; x < GridSize; x++ {
		if !columnEmpty(g, x, posY, endY) {
			break
		}
	}
	x -= 1
	return x
}

func columnEmpty(g *gameGrid, x, y, endY int) bool {
	for ; y <= endY; y++ {
		if !g.isEmpty(x, y) {
			return false
		}
	}
	return true
}
