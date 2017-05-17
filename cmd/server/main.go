package main

import (
	"io"
	"log"
	"net"

	"github.com/flazz/grpctalk/point"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	maxX = 10
	maxY = 10
)

func main() {
	var gs gameserver

	server := grpc.NewServer()
	point.RegisterGameServer(server, &gs)
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatal("error listening:", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatal("error serving: ", err)
	}

}

type gameserver struct{}

func (gs *gameserver) SetGoal(ctx context.Context, req *point.SetGoalRequest) (*point.SetGoalResponse, error) {
	// TODO do some work
	log.Println("SET-GOAL", req)
	return &point.SetGoalResponse{}, nil
}

func (gs *gameserver) Move(clientStream point.Game_MoveServer) error {
	log.Println("MOVE")
	for {
		req, err := clientStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return grpc.Errorf(codes.Internal, "%v", err)
		}
		log.Println("MOVE", req)
	}

	return nil
}

func (gs *gameserver) Watch(req *point.WatchRequest, serverStream point.Game_WatchServer) error {
	log.Println("WATCH", req)
	for i := 0; i < 3; i++ {
		log.Println("WATCH SEND")
		err := serverStream.Send(&point.WatchResponse{
			Name:     req.Name,
			Position: &point.Point{X: int32(i), Y: int32(i)},
			Score:    false,
		})
		if err != nil {
			return grpc.Errorf(codes.Internal, "%v", err)
		}

	}

	return nil
}

func (gs *gameserver) Chat(biStream point.Game_ChatServer) error {
	log.Println("CHAT")

	for {
		req, err := biStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return grpc.Errorf(codes.Internal, "%v", err)
		}
		log.Println("CHAT-IN", req)

		if err := biStream.Send(&point.ChatResponse{Msg: req.Msg, Name: req.Name}); err != nil {
			return grpc.Errorf(codes.Internal, "%v", err)
		}
		log.Println("CHAT-OUT")

	}

}
