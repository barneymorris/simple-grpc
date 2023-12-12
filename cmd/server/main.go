package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/barneymorris/simple-grpc/internal/config"
	"github.com/barneymorris/simple-grpc/pkg/note_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	note_v1.UnimplementedNoteV1Server
	pool *pgxpool.Pool
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
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %s", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %s", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen grpc server: %s", err)
	}

	pgContext, cancel := context.WithTimeout(context.Background(), 5000)
	pool, err := pgxpool.New(pgContext, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer cancel()
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	note_v1.RegisterNoteV1Server(s, server{pool: pool})

	log.Printf("server listening at %s", grpcConfig.Address())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}