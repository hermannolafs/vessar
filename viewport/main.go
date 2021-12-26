package main

import (
	"fmt"
	"github.com/kelindar/tile"
	"github.com/logrusorgru/aurora/v3"
	"github.com/nsf/termbox-go"
	"hermannolafs/vessar/viewport/input"
	"log"
	"strconv"
)

var (
	zeroZeroTile = tile.Tile{255, 255, 255, 0, 0, 0}
	maxMaxTile   = tile.Tile{0, 0, 0, 255, 255, 255}
	middletile   = tile.Tile{64, 64, 64, 128, 128, 128}

	standardTile = tile.Tile{9, 9, 9, 3, 3, 3}
	equalPosTile = tile.Tile{21, 21, 21, 9, 9, 9}
)

const (
	defaultPlayerViewSize = 5
)

func main() {
	grid := new21Grid()

	log.Print("Created Grid of size ", aurora.Green(grid.Size))

	playerView := setupPlayerView(grid)

	log.Printf("Initializing Termbox")
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	defer termbox.Close()
	log.Printf("Termbox initialized")

	log.Printf("Initializing input listener")
	playerView.input.Start()
	defer playerView.input.Stop()
	log.Printf("Input listener initialized")

mainloop:
	for {
		select {
		case ev := <-playerView.input.Queue:
			if ev.Key == termbox.KeyCtrlC {
				println("GOT CTRL+C")
				break mainloop
			}
		}
		playerView.view.Each(printTile)
	}
}

func printView(view *tile.View) {
	view.Each(printTile)
	_ = termbox.Flush()
}

func printTile(point tile.Point, t tile.Tile) {
	termbox.SetCell(
		int(point.X),
		int(point.Y),
		'0',
		termbox.ColorBlue,
		termbox.ColorRed,
	)
}



func pointToInts(point tile.Point) (int, int) {
	return int(point.X), int(point.Y)
}

func setupPlayerView(grid *tile.Grid) PlayerView {
	setReferenceTiles(grid)
	return newPlayerView(grid)
}

type PlayerView struct {
	view *tile.View
	size tile.Point // represents grid size of view
	input *input.Input

}

func newPlayerView(grid *tile.Grid) PlayerView {
	sizePos := tile.Point{defaultPlayerViewSize, defaultPlayerViewSize}

	playerView := PlayerView{
		view: grid.View(
			rectFromTwoPositions(
				tile.Point{0, 0},
				sizePos,
			),
			func(p tile.Point, tile tile.Tile) {},
		),
		size: sizePos, // Maybe this should be a function instead of private field?
		input: input.NewInput(),
	}
	return playerView
}
// (x1,y1), (x2,y2) --> rect(x1,y1,x2,y2)
// Wraps tile.New:w
func rectFromTwoPositions(lowPosition, highPosition tile.Point) tile.Rect {
	return tile.NewRect(
		lowPosition.X, lowPosition.Y,
		highPosition.X, highPosition.Y,
	)
}

func printGrid(grid *tile.Grid) {
	for y := int16(0); y < grid.Size.Y; y++ {
		for x := int16(0); x < grid.Size.X; x++ {
			currentTile, _ := grid.At(x, y)
			fg, bg := getColorFromTileAvg(currentTile)
			fmt.Print(aurora.Index(fg, aurora.BgIndex(bg, posToString(x, y))))
		}
		println()
	}
}

// x=1, y=1 --> "1,1 "
func posToString(x int16, y int16) string {
	return strconv.Itoa(int(x)) + "," +
		strconv.Itoa(int(y)) + " "
}

// Calculates bg color from vales in [0:2] and fg from [3:5]
// Used temporarily as placeholder for proper ASCII/Sprites
func getColorFromTileAvg(currentTile tile.Tile) (uint8, uint8) {
	fgAverage := (currentTile[0] + currentTile[1] + currentTile[2]) / 3
	bgAverage := (currentTile[3] + currentTile[4] + currentTile[5]) / 3

	if fgAverage == bgAverage && fgAverage == 0 {
		return 0, 23
	}

	return fgAverage, bgAverage
}

// Sets 00, X,Y and X/2,Y/2 as standard tiles
func setReferenceTiles(grid *tile.Grid) {
	grid.Each(func(point tile.Point, _ tile.Tile) {
		grid.WriteAt(point.X, point.Y, standardTile)
	})

	// Marking axis
	for i := int16(0); i < grid.Size.X; i++ {
		grid.WriteAt(i, i, equalPosTile)
	}

	// Mark Corners for reference
	grid.WriteAt(0, 0, zeroZeroTile)
	grid.WriteAt(grid.Size.X-1, grid.Size.Y-1, maxMaxTile)
	// Mark middle
	grid.WriteAt(grid.Size.X/2, grid.Size.Y/2, middletile)

	//zeroTile, _ := grid.At(0,0)
	//log.Print("Set reference tiles, 0,0 is now: ", aurora.Sprintf(aurora.Blue("%v"), zeroTile))
}

// Returns 9x9 grid
func new9Grid() *tile.Grid {
	return tile.NewGrid(9, 9)
}

// Returns 21x21 grid
func new21Grid() *tile.Grid {
	return tile.NewGrid(21, 21)
}