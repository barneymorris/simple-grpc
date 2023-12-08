package main

import (
	"context"
	"log"
	"time"

	"github.com/barneymorris/simple-grpc/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	noteID = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to grpc server: %s", err)
	}

	defer conn.Close()

	c := note_v1.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &note_v1.GetRequest{Id: noteID})
	if err != nil {
		log.Fatalf("failed to get note with id %d, error: %s", noteID, err)
	}

	log.Printf("Note info: %+v", r.GetNote())
}