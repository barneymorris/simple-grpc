package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/barneymorris/simple-grpc/pkg/note_v1"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	note_v1.UnimplementedNoteV1Server
}

func (s server) Get(ctx context.Context, req *note_v1.GetRequest) (*note_v1.GetResponse, error) {
	log.Printf("Node id: %d", req.Id)

	return &note_v1.GetResponse{
		Note: &note_v1.Note{
			Id: req.GetId(),
			Info: &note_v1.NoteInfo{
				Title: gofakeit.BeerName(),
				Content: gofakeit.IPv4Address(),
				Author: gofakeit.Name(),
				IsPublic: gofakeit.Bool(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	impl := server{}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen grpc server: %s", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	note_v1.RegisterNoteV1Server(s, impl)

	log.Printf("server listening at port %d", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}