package main

import (
	"log"

	"github.com/flazz/grpctalk/point"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// make a client
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("dial", err)
	}
	defer conn.Close()
	client := point.NewGameClient(conn)

	// unary
	uRep, err := client.SetGoal(context.Background(), &point.SetGoalRequest{
		Position: &point.Point{X: 10, Y: 20},
	})
	log.Println("unary", uRep, err)

	// server stream
	serverStream, err := client.Watch(context.Background(), &point.WatchRequest{Name: "foo"})
	log.Println("serverStream err", err)
	for {
		ssRep, err := serverStream.Recv()
		log.Println("serverStream", ssRep, err)
		if err != nil {
			break
		}
	}

	// client stream
	clientStream, err := client.Move(context.Background())
	log.Println("clientStream err", err)
	for i := 0; i < 3; i++ {
		err := clientStream.Send(&point.MoveRequest{})
		log.Println("clientStream err", err)
	}
	scRep, err := clientStream.CloseAndRecv()
	log.Println("clientStream", scRep, err)

	// bi-directional strea
	biStream, err := client.Chat(context.Background())
	for i := 0; i < 3; i++ {
		err := biStream.Send(&point.ChatRequest{Msg: "the msg", Name: "the name"})
		log.Println("biStream.Send err", err)

		bsRep, err := biStream.Recv()
		log.Println("biStream.Recv", bsRep, err)
	}
	err = biStream.CloseSend()
	log.Println("biStream.CloseSend", err)
}
