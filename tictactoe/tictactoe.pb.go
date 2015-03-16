// Code generated by protoc-gen-go.
// source: tictactoe.proto
// DO NOT EDIT!

/*
Package tictactoe is a generated protocol buffer package.

It is generated from these files:
	tictactoe.proto

It has these top-level messages:
	CreateRequest
	CreateReply
	TurnRequest
	TurnReply
	Event
*/
package tictactoe

import proto "github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/github.com/golang/protobuf/proto"

import (
	context "github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/golang.org/x/net/context"
	grpc "github.com/protogalaxy/service-tictactoe-game/Godeps/_workspace/src/google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type Mark int32

const (
	Mark_EMPTY Mark = 0
	Mark_X     Mark = 1
	Mark_Y     Mark = 2
)

var Mark_name = map[int32]string{
	0: "EMPTY",
	1: "X",
	2: "Y",
}
var Mark_value = map[string]int32{
	"EMPTY": 0,
	"X":     1,
	"Y":     2,
}

func (x Mark) String() string {
	return proto.EnumName(Mark_name, int32(x))
}

type CreateReply_ResponseStatus int32

const (
	CreateReply_SUCCESS CreateReply_ResponseStatus = 0
)

var CreateReply_ResponseStatus_name = map[int32]string{
	0: "SUCCESS",
}
var CreateReply_ResponseStatus_value = map[string]int32{
	"SUCCESS": 0,
}

func (x CreateReply_ResponseStatus) String() string {
	return proto.EnumName(CreateReply_ResponseStatus_name, int32(x))
}

type TurnReply_ResponseStatus int32

const (
	TurnReply_SUCCESS           TurnReply_ResponseStatus = 0
	TurnReply_INVALID_MOVE      TurnReply_ResponseStatus = 1
	TurnReply_NOT_ACTIVE_PLAYER TurnReply_ResponseStatus = 2
	TurnReply_FINISHED          TurnReply_ResponseStatus = 3
)

var TurnReply_ResponseStatus_name = map[int32]string{
	0: "SUCCESS",
	1: "INVALID_MOVE",
	2: "NOT_ACTIVE_PLAYER",
	3: "FINISHED",
}
var TurnReply_ResponseStatus_value = map[string]int32{
	"SUCCESS":           0,
	"INVALID_MOVE":      1,
	"NOT_ACTIVE_PLAYER": 2,
	"FINISHED":          3,
}

func (x TurnReply_ResponseStatus) String() string {
	return proto.EnumName(TurnReply_ResponseStatus_name, int32(x))
}

type Event_Type int32

const (
	Event_GAME_CREATED Event_Type = 0
	Event_TURN_PLAYED  Event_Type = 1
)

var Event_Type_name = map[int32]string{
	0: "GAME_CREATED",
	1: "TURN_PLAYED",
}
var Event_Type_value = map[string]int32{
	"GAME_CREATED": 0,
	"TURN_PLAYED":  1,
}

func (x Event_Type) String() string {
	return proto.EnumName(Event_Type_name, int32(x))
}

type CreateRequest struct {
	UserIds []string `protobuf:"bytes,1,rep,name=user_ids" json:"user_ids,omitempty"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}

type CreateReply struct {
	Status CreateReply_ResponseStatus `protobuf:"varint,1,opt,name=status,enum=tictactoe.CreateReply_ResponseStatus" json:"status,omitempty"`
	GameId string                     `protobuf:"bytes,2,opt,name=game_id" json:"game_id,omitempty"`
}

func (m *CreateReply) Reset()         { *m = CreateReply{} }
func (m *CreateReply) String() string { return proto.CompactTextString(m) }
func (*CreateReply) ProtoMessage()    {}

type TurnRequest struct {
	GameId string            `protobuf:"bytes,1,opt,name=game_id" json:"game_id,omitempty"`
	UserId string            `protobuf:"bytes,2,opt,name=user_id" json:"user_id,omitempty"`
	Move   *TurnRequest_Move `protobuf:"bytes,3,opt,name=move" json:"move,omitempty"`
}

func (m *TurnRequest) Reset()         { *m = TurnRequest{} }
func (m *TurnRequest) String() string { return proto.CompactTextString(m) }
func (*TurnRequest) ProtoMessage()    {}

func (m *TurnRequest) GetMove() *TurnRequest_Move {
	if m != nil {
		return m.Move
	}
	return nil
}

type TurnRequest_Move struct {
	X int32 `protobuf:"varint,1,opt,name=x" json:"x,omitempty"`
	Y int32 `protobuf:"varint,2,opt,name=y" json:"y,omitempty"`
}

func (m *TurnRequest_Move) Reset()         { *m = TurnRequest_Move{} }
func (m *TurnRequest_Move) String() string { return proto.CompactTextString(m) }
func (*TurnRequest_Move) ProtoMessage()    {}

type TurnReply struct {
	Status TurnReply_ResponseStatus `protobuf:"varint,1,opt,name=status,enum=tictactoe.TurnReply_ResponseStatus" json:"status,omitempty"`
	Winner *TurnReply_Winner        `protobuf:"bytes,2,opt,name=winner" json:"winner,omitempty"`
}

func (m *TurnReply) Reset()         { *m = TurnReply{} }
func (m *TurnReply) String() string { return proto.CompactTextString(m) }
func (*TurnReply) ProtoMessage()    {}

func (m *TurnReply) GetWinner() *TurnReply_Winner {
	if m != nil {
		return m.Winner
	}
	return nil
}

type TurnReply_Winner struct {
	Draw   bool   `protobuf:"varint,1,opt,name=draw" json:"draw,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id" json:"user_id,omitempty"`
}

func (m *TurnReply_Winner) Reset()         { *m = TurnReply_Winner{} }
func (m *TurnReply_Winner) String() string { return proto.CompactTextString(m) }
func (*TurnReply_Winner) ProtoMessage()    {}

type Event struct {
	Type       Event_Type               `protobuf:"varint,1,opt,name=type,enum=tictactoe.Event_Type" json:"type,omitempty"`
	Timestamp  *Event_Timestamp         `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	GameId     string                   `protobuf:"bytes,3,opt,name=game_id" json:"game_id,omitempty"`
	UserId     string                   `protobuf:"bytes,4,opt,name=user_id" json:"user_id,omitempty"`
	UserList   []string                 `protobuf:"bytes,5,rep,name=user_list" json:"user_list,omitempty"`
	Move       *TurnRequest_Move        `protobuf:"bytes,6,opt,name=move" json:"move,omitempty"`
	TurnStatus TurnReply_ResponseStatus `protobuf:"varint,7,opt,name=turn_status,enum=tictactoe.TurnReply_ResponseStatus" json:"turn_status,omitempty"`
	Winner     *TurnReply_Winner        `protobuf:"bytes,8,opt,name=winner" json:"winner,omitempty"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}

func (m *Event) GetTimestamp() *Event_Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Event) GetMove() *TurnRequest_Move {
	if m != nil {
		return m.Move
	}
	return nil
}

func (m *Event) GetWinner() *TurnReply_Winner {
	if m != nil {
		return m.Winner
	}
	return nil
}

type Event_Timestamp struct {
	Seconds int64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
	Nanos   int32 `protobuf:"varint,2,opt,name=nanos" json:"nanos,omitempty"`
}

func (m *Event_Timestamp) Reset()         { *m = Event_Timestamp{} }
func (m *Event_Timestamp) String() string { return proto.CompactTextString(m) }
func (*Event_Timestamp) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("tictactoe.Mark", Mark_name, Mark_value)
	proto.RegisterEnum("tictactoe.CreateReply_ResponseStatus", CreateReply_ResponseStatus_name, CreateReply_ResponseStatus_value)
	proto.RegisterEnum("tictactoe.TurnReply_ResponseStatus", TurnReply_ResponseStatus_name, TurnReply_ResponseStatus_value)
	proto.RegisterEnum("tictactoe.Event_Type", Event_Type_name, Event_Type_value)
}

// Client API for GameManager service

type GameManagerClient interface {
	CreateGame(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error)
	PlayTurn(ctx context.Context, in *TurnRequest, opts ...grpc.CallOption) (*TurnReply, error)
}

type gameManagerClient struct {
	cc *grpc.ClientConn
}

func NewGameManagerClient(cc *grpc.ClientConn) GameManagerClient {
	return &gameManagerClient{cc}
}

func (c *gameManagerClient) CreateGame(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error) {
	out := new(CreateReply)
	err := grpc.Invoke(ctx, "/tictactoe.GameManager/CreateGame", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameManagerClient) PlayTurn(ctx context.Context, in *TurnRequest, opts ...grpc.CallOption) (*TurnReply, error) {
	out := new(TurnReply)
	err := grpc.Invoke(ctx, "/tictactoe.GameManager/PlayTurn", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GameManager service

type GameManagerServer interface {
	CreateGame(context.Context, *CreateRequest) (*CreateReply, error)
	PlayTurn(context.Context, *TurnRequest) (*TurnReply, error)
}

func RegisterGameManagerServer(s *grpc.Server, srv GameManagerServer) {
	s.RegisterService(&_GameManager_serviceDesc, srv)
}

func _GameManager_CreateGame_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(CreateRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(GameManagerServer).CreateGame(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _GameManager_PlayTurn_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(TurnRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(GameManagerServer).PlayTurn(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _GameManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tictactoe.GameManager",
	HandlerType: (*GameManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGame",
			Handler:    _GameManager_CreateGame_Handler,
		},
		{
			MethodName: "PlayTurn",
			Handler:    _GameManager_PlayTurn_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
