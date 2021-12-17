package main

import (
	"bytes"
	"context"
	"github.com/kelindar/tile"
	"github.com/logrusorgru/aurora/v3"
	"google.golang.org/grpc"
	"hermannolafs/vessar/beinagrind"
	samskipti "hermannolafs/vessar/samskipti/go"
	"log"
	"strconv"
	"time"
)

type Skermur struct {
	kort samskipti.KortClient
}

func main() {
	skermur := NewSkermur()

	skermur.Start()
}

func NewSkermur() Skermur {
	skermur := Skermur{
		kort: connectKort(),
	}

	return skermur
}

func (skermur Skermur) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	kort, err := skermur.kort.GetKort(ctx, &samskipti.Null{})
	if err != nil {
		log.Fatalf("could not fetch Kort : ", err)
	}

	log.Print(
		aurora.Blue("SKERMUR | "),
		"Got tiles for Kort:",
		aurora.Green(kort.GetTiles()),
	)

	grid, err := tile.ReadFrom(bytes.NewReader(kort.GetTiles()))
	if err != nil {
		panic(err)
	}

	_ = grid
}

func connectKort() samskipti.KortClient {
	log.Printf("connecting to Kort")
	kortConnection, err := grpc.Dial("localhost:"+strconv.Itoa(int(beinagrind.KortPort)), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	return samskipti.NewKortClient(kortConnection)
}
