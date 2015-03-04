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

//go:generate protoc --go_out=plugins=grpc:. -I ../protos ../protos/tictactoe.proto

package tictactoe

import (
	"errors"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"golang.org/x/net/context"
)

const GridSize int = 3

type gameGrid struct {
	grid []Mark
}

func (g *gameGrid) set(x, y int, m Mark) {
	g.grid[y*GridSize+x] = m
}

func newGameGrid() *gameGrid {
	return &gameGrid{
		grid: make([]Mark, GridSize*GridSize),
	}
}

type GameID string

type game struct {
	ID      GameID
	Grid    *gameGrid
	Players map[string]Mark
}

func newGame(ID GameID, playerOne, playerTwo string) *game {
	return &game{
		ID:   ID,
		Grid: newGameGrid(),
		Players: map[string]Mark{
			playerOne: Mark_X,
			playerTwo: Mark_Y,
		},
	}
}

var (
	ErrInvalidMove = errors.New("invalid move")
)

func (g *game) placeMark(userID string, x, y int) error {
	g.Grid.set(x, y, g.Players[userID])
	return nil
}

type GameManager struct {
	lock        sync.Mutex
	activeGames map[GameID]*game
}

func NewGameManager() *GameManager {
	return &GameManager{
		activeGames: make(map[GameID]*game),
	}
}

func newID() GameID {
	return GameID(uuid.NewRandom().String())
}

func (m *GameManager) CreateGame(ctx context.Context, req *CreateRequest) (*CreateReply, error) {
	if len(req.UserIds) != 2 {
		return nil, errors.New("number of players must be 2")
	}

	var rep CreateReply

	game := newGame(newID(), req.UserIds[0], req.UserIds[1])

	m.lock.Lock()
	m.activeGames[game.ID] = game
	m.lock.Unlock()

	rep.Status = ResponseStatus_SUCCESS
	rep.GameId = string(game.ID)
	return &rep, nil
}

func (m *GameManager) PlayTurn(ctx context.Context, req *TurnRequest) (*TurnReply, error) {
	var rep TurnReply

	m.lock.Lock()
	defer m.lock.Unlock()

	game := m.activeGames[GameID(req.GameId)]
	game.placeMark(req.UserId, int(req.Move.X), int(req.Move.Y))

	rep.Status = TurnReply_SUCCESS
	return &rep, nil
}