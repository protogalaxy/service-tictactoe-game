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
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/Shopify/sarama"
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/golang/protobuf/proto"
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/golang.org/x/net/context"
)

const streamTopic = "tictactoe-game-events"

type GameManager struct {
	lock        sync.Mutex
	activeGames map[GameID]*game

	stream sarama.SyncProducer
}

func NewGameManager(s sarama.SyncProducer) *GameManager {
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
	if err := sendMessage(m.stream, streamTopic, ev.GameId, &ev); err != nil {
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
		rep.Status = TurnReply_FINISHED
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

	rep.MoveId = game.lastMoveID()

	ev := Event{
		Type:      Event_TURN_PLAYED,
		Timestamp: time.Now().UnixNano(),
		GameId:    string(req.GameId),
		UserId:    req.UserId,
		UserList:  game.PlayerList,

		Move:       req.Move,
		TurnStatus: rep.Status,
		MoveId:     rep.MoveId,

		NextPlayer: game.activePlayer(),
	}
	if game.isFinished() {
		rep.Status = TurnReply_FINISHED
		ev.Winner = game.winner()
	} else {
		ev.ValidMoves = game.validMoves()
	}
	if err := sendMessage(m.stream, streamTopic, ev.GameId, &ev); err != nil {
		return nil, err
	}

	return &rep, nil
}

func sendMessage(p sarama.SyncProducer, topic string, key string, m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		glog.Fatalf("Encoding message: %s", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
		Key:   sarama.StringEncoder(key),
	}
	_, _, err = p.SendMessage(msg)
	return err
}
