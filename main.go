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
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/Shopify/sarama"
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/google.golang.org/grpc"
	"github.com/protogalaxy/service-tictactoe-game/tictactoe"
)

var port = flag.Int("port", 9090, "port to listen on")

func ParseLinkEnv(name string) string {
	v := os.Getenv(name + "_PORT")
	if v == "" {
		glog.Fatalf("Missing environment variable %s_PORT", name)
	}
	connStr := strings.Split(v, "//")[1]
	return connStr
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	cfg := sarama.NewConfig()
	cfg.ClientID = "service-tictactoe-game"
	producer, err := sarama.NewSyncProducer([]string{ParseLinkEnv("KAFKA")}, cfg)
	if err != nil {
		glog.Fatalf("Unable to connect to kafka: %s", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			glog.Errorf("Error closing producer: %s", err)
		}
	}()

	grpcServer := grpc.NewServer()
	tictactoe.RegisterGameManagerServer(grpcServer, tictactoe.NewGameManager(producer))
	grpcServer.Serve(socket)
}
