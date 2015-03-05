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

func (g *gameGrid) checkDiagonals() bool {
	return g.checkDiagonalDown()
}

func (g *gameGrid) checkDiagonalDown() bool {
	first := g.get(0, 0)
	return first == g.get(1, 1) && first == g.get(2, 2)
}

func (g *gameGrid) checkDiagonalUp() bool {
	first := g.get(0, 2)
	return first == g.get(1, 1) && first == g.get(2, 0)
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

func (g *gameGrid) coordinatesValid(x, y int) bool {
	return validateIndex(x) && validateIndex(y)
}

func validateIndex(x int) bool {
	return x >= 0 && x < GridSize
}

type GameID string

type game struct {
	ID            GameID
	Grid          *gameGrid
	CurrentPlayer int
	PlayerList    []string
	Players       map[string]Mark
	GameFinished  bool
	Winner        string
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
)

func (g *game) activePlayer() string {
	return g.PlayerList[g.CurrentPlayer]
}

func (g *game) updateActivePlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.PlayerList)
}

func (g *game) checkWinner(userID string, x, y int) {
	if g.Grid.checkMarksHorizontal(y) || g.Grid.checkMarksVertical(x) || g.Grid.checkDiagonals() {
		g.Winner = userID
		g.GameFinished = true
	} else if g.Grid.isFull() {
		g.GameFinished = true
	}
}

func (g *game) isFinished() bool {
	return g.GameFinished
}

func (g *game) isDraw() bool {
	return g.GameFinished && g.Winner == ""
}

func (g *game) winner() string {
	return g.Winner
}

func (g *game) placeMark(userID string, x, y int) error {
	if g.activePlayer() != userID {
		return ErrNotActivePlayer
	} else if !g.Grid.coordinatesValid(x, y) || !g.Grid.isEmpty(x, y) {
		return ErrInvalidMove
	}
	g.Grid.set(x, y, g.Players[userID])
	g.checkWinner(userID, x, y)
	g.updateActivePlayer()
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
	rep.Status = TurnReply_SUCCESS

	m.lock.Lock()
	defer m.lock.Unlock()

	game := m.activeGames[GameID(req.GameId)]
	if game.isFinished() {
		prepareWinnerResponse(&rep, game)
		return &rep, nil
	}

	err := game.placeMark(req.UserId, int(req.Move.X), int(req.Move.Y))
	switch {
	case err == ErrInvalidMove:
		rep.Status = TurnReply_INVALID_MOVE
	case err == ErrNotActivePlayer:
		rep.Status = TurnReply_NOT_ACTIVE_PLAYER
	case err != nil:
		return nil, err
	}

	if game.isFinished() {
		prepareWinnerResponse(&rep, game)
	}

	return &rep, nil
}

func prepareWinnerResponse(rep *TurnReply, game *game) {
	rep.Status = TurnReply_FINISHED
	rep.Winner = &TurnReply_Winner{
		Draw:   game.isDraw(),
		UserId: game.winner(),
	}
}
