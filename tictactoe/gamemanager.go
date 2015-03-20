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
	"time"

	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid"
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/protogalaxy/service-tictactoe-game/stream"
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

type GameManager struct {
	lock        sync.Mutex
	activeGames map[GameID]*game

	stream stream.ProtoProducer
}

func NewGameManager(s stream.ProtoProducer) *GameManager {
	return &GameManager{
		activeGames: make(map[GameID]*game),
		stream:      s,
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

	gameID := newID()
	game := newGame(gameID, req.UserIds[0], req.UserIds[1])

	m.lock.Lock()
	m.activeGames[game.ID] = game
	m.lock.Unlock()

	rep.Status = CreateReply_SUCCESS
	rep.GameId = string(game.ID)

	ev := Event{
		Type:      Event_GAME_CREATED,
		Timestamp: time.Now().UnixNano(),
		GameId:    string(gameID),
		UserId:    req.UserIds[0],
		UserList:  req.UserIds,

		NextPlayer: game.activePlayer(),
		ValidMoves: game.validMoves(),
	}
	if err := m.stream.SendMessage(&ev); err != nil {
		return nil, err
	}

	return &rep, nil
}

func (m *GameManager) PlayTurn(ctx context.Context, req *TurnRequest) (*TurnReply, error) {
	var rep TurnReply
	rep.Status = TurnReply_SUCCESS

	m.lock.Lock()
	defer m.lock.Unlock()

	game := m.activeGames[GameID(req.GameId)]
	if game.isFinished() {
		rep.Winner = game.winner()
		return &rep, nil
	}

	err := game.placeMark(req.UserId, req.MoveId, int(req.Move.X), int(req.Move.Y))
	switch {
	case err == ErrInvalidMove:
		rep.Status = TurnReply_INVALID_MOVE
	case err == ErrNotActivePlayer:
		rep.Status = TurnReply_NOT_ACTIVE_PLAYER
	case err == ErrInvalidMoveID:
		rep.Status = TurnReply_INVALID_MOVE_ID
	case err != nil:
		return nil, err
	}

	if game.isFinished() {
		rep.Winner = game.winner()
	}

	ev := Event{
		Type:      Event_TURN_PLAYED,
		Timestamp: time.Now().UnixNano(),
		GameId:    string(req.GameId),
		UserId:    req.UserId,
		UserList:  game.PlayerList,

		Move:       req.Move,
		TurnStatus: rep.Status,
		MoveId:     game.lastMoveID(),

		NextPlayer: game.activePlayer(),
	}
	if rep.Winner != nil {
		ev.Winner = rep.Winner
	} else {
		ev.ValidMoves = game.validMoves()
	}
	if err := m.stream.SendMessage(&ev); err != nil {
		return nil, err
	}

	return &rep, nil
}
