package main

import (
	"context"
	"github.com/logrusorgru/aurora/v3"
	"google.golang.org/grpc"
	"hermannolafs/vessar/beinagrind"
	samskipti "hermannolafs/vessar/samskipti/go"
	"log"
	"strconv"
	"time"
)

type Hjarta struct {
	heili         samskipti.HeiliClient
	ModelMetadata samskipti.ModelMetadata
}

func main() {
	hjarta := NewHjarta()

	hjarta.Start()
}

func NewHjarta() *Hjarta {
	return &Hjarta{
		heili: connectHeili(),
		ModelMetadata: samskipti.ModelMetadata{
			Port:      beinagrind.HjartaPort,
			ModelName: 1,
		},
	}
}

func (hjarta Hjarta) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	registration, err := hjarta.heili.Register(ctx, &hjarta.ModelMetadata)
	if err != nil {
		log.Fatalf("could not register to Heili: ", err)
	}

	log.Print(
		aurora.Blue(hjarta.ModelMetadata.ModelName),
		"registering humour was:",
		aurora.Green(registration.GetHumourType()),
	)
}

func connectHeili() samskipti.HeiliClient {
	log.Printf("connecting heili")
	heilaConnection, err := grpc.Dial("localhost:"+strconv.Itoa(int(beinagrind.HeiliPort)), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	return samskipti.NewHeiliClient(heilaConnection)
}
