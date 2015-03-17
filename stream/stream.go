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

package stream

import (
	"fmt"

	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/golang/protobuf/proto"
)

type Producer interface {
	Send(b []byte) error
}

type ProtoProducer interface {
	SendMessage(m proto.Message) error
	Producer
}

type TracingProducer struct {
}

func NewProducer() *TracingProducer {
	return &TracingProducer{}
}

func (s *TracingProducer) Send(b []byte) error {
	fmt.Println("Event size: ", len(b))
	return nil
}

func (s *TracingProducer) SendMessage(m proto.Message) error {
	fmt.Println("Evcent: ", m)
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return s.Send(b)
}
