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

package main

import (
	"flag"
	"math/rand"
	"net"
	"time"

	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/google.golang.org/grpc"
	"github.com/protogalaxy/service-tictactoe-game/stream"
	"github.com/protogalaxy/service-tictactoe-game/tictactoe"
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	socket, err := net.Listen("tcp", ":9090")
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	producer := stream.NewProducer()

	grpcServer := grpc.NewServer()
	tictactoe.RegisterGameManagerServer(grpcServer, tictactoe.NewGameManager(producer))
	grpcServer.Serve(socket)
}
