package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	samskipti "hermannolafs/vessar/samskipti/go"

	"log"
	"net"
	"strconv"
	"sync"
)

func main() {
	heili := NewHeili()

	heili.Start()

}

// Four lobes?
type Heili struct {
	líffæri map[samskipti.HumorType][]samskipti.HjartaClient
	lock    sync.RWMutex
	port    int
}

func NewHeili() *Heili {
	heili := Heili{
		líffæri: make(map[samskipti.HumorType][]samskipti.HjartaClient),
		lock:    sync.RWMutex{},
		port:    4000,
	}

	return &heili
}

func (heili *Heili) Start() {
	log.Printf("Starting Heili")

	listener, err := net.Listen(
		"tcp",
		":"+strconv.Itoa(heili.port),
	)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	samskipti.RegisterHeiliServer(server, heili)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}

func (heili *Heili) Register(ctx context.Context, metadata *samskipti.ModelMetadata) (*samskipti.Humor, error) {
	heili.lock.Lock()
	defer heili.lock.Unlock()
	log.Printf("Registering new Hjarta")

	hjartaConnection, err := grpc.Dial("localhost:"+string(metadata.Port), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	hjartaClient := samskipti.NewHjartaClient(hjartaConnection)

	// if there is no list yet, create one
	if _, ok := heili.líffæri[metadata.HumourType]; !ok {
		heili.líffæri[metadata.HumourType] = make([]samskipti.HjartaClient, 0)
	}

	// append to slice
	heili.líffæri[metadata.HumourType] = append(heili.líffæri[metadata.HumourType], hjartaClient)

	return &samskipti.Humor{
		HumourType: metadata.HumourType,
	}, nil
}

func (heili *Heili) Deregister(ctx context.Context, metadata *samskipti.ModelMetadata) (*samskipti.Humor, error) {
	panic("implement me")
}
