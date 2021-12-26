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
	skermur := Skermur{}

	return skermur
}

func (skermur *Skermur) Start() {
	if err := skermur.connectKort(); err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 3)

		if _, err := skermur.GetKort(); err != nil {
			log.Print(aurora.BgRed("|ERROR|"), " Could not run GetKort")
			log.Print(aurora.BgRed(err.Error()))
		}
	}
}


func (skermur *Skermur) GetKort() (*tile.Grid, error) {
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

	return tile.ReadFrom(bytes.NewReader(kort.GetTiles()))
}

// TODO timeout here or in Start()?
func (skermur *Skermur) connectKort() error {
	log.Printf("connecting to Kort")

	kortURL := "localhost:"+strconv.Itoa(int(beinagrind.KortPort))
	kortConnection, err := grpc.Dial(kortURL, grpc.WithInsecure())

	if err != nil {
		log.Print(aurora.BgRed("Could not connect to kort on URL "), aurora.Green(kortURL))
		return err
	}

	skermur.kort = samskipti.NewKortClient(kortConnection)
	return err
}
