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

func (heili *Heili) Register(ctx context.Context, metadata *samskipti.ModelMetadata) (*samskipti.ModelMetadata, error) {
	heili.lock.Lock()
	defer heili.lock.Unlock()
	log.Printf("Registering new Hjarta")

	hjartaClient := connectHjartaOnLocalhost(metadata)

	// if there is no list yet, create one

	heili.register(metadata, hjartaClient)

	log.Printf("Current list: \n %+v\n")

	return metadata, nil
}

func (heili *Heili)register(metadata *samskipti.ModelMetadata, hjartaClient samskipti.HjartaClient) {
	for _, humour := range metadata.GetHumourList() {
		log.Printf("Cast it to: %v", humour.GetHumour())

		if _, ok := heili.líffæri[humour.GetHumour()]; !ok {
			heili.líffæri[humour.GetHumour()] = make([]samskipti.HjartaClient, 0)
		}
		heili.líffæri[humour.GetHumour()] = append(heili.líffæri[humour.GetHumour()], hjartaClient)
	}
}

func connectHjartaOnLocalhost(metadata *samskipti.ModelMetadata) samskipti.HjartaClient {
	hjartaConnection, err := grpc.Dial("localhost:"+string(metadata.Port), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	return samskipti.NewHjartaClient(hjartaConnection)
}

func (heili *Heili) Deregister(ctx context.Context, metadata *samskipti.ModelMetadata) (*samskipti.ModelMetadata, error) {
	
}

func (heili *Heili) GetRegisteredModels(ctx context.Context, null *samskipti.Null) (*samskipti.ModelMetadataList, error) {
	panic("implement me")
}



