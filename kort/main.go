package main

import (
	"bytes"
	"context"
	"github.com/kelindar/tile"
	"github.com/logrusorgru/aurora/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hermannolafs/vessar/beinagrind"
	samskipti "hermannolafs/vessar/samskipti/go"
	"log"
	"net"
	"strconv"
)

type Kort struct {
	reitir *tile.Grid
	port   int
}

func (kort *Kort) Start() {
	log.Printf("Starting Kort")

	listener, err := net.Listen(
		"tcp",
		":"+strconv.Itoa(kort.port),
	)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	samskipti.RegisterKortServer(server, kort)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}

func createNewTestGrid() *tile.Grid {
	grid := tile.NewGrid(3, 3)
	grid.WriteAt(1, 1, tile.Tile{1, 2, 3, 4, 5, 6})

	return grid
}

func createNewTestKort() *Kort {
	newKort := Kort{
		reitir: createNewTestGrid(),
		port:   int(beinagrind.KortPort),
	}

	return &newKort
}

func main() {
	kort := createNewTestKort()

	kort.Start()
}

func (kort *Kort) GetKort(ctx context.Context, null *samskipti.Null) (*samskipti.KortBytes, error) {
	log.Printf("Getting kort")
	output := new(bytes.Buffer)

	if _, err := kort.reitir.WriteTo(output); err != nil {
		return nil, err
	}

	log.Print(aurora.Blue("GetKort"),
		"Got this output written: ",
		aurora.Green(output.Bytes()))

	return &samskipti.KortBytes{
		Tiles: output.Bytes(),
	}, nil
}
