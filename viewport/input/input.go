package input

import (
	"github.com/nsf/termbox-go"
)

type Input struct {
	Queue chan termbox.Event
	EndKey termbox.Key
	Ctrl   chan bool
}

func NewInput() *Input {
	return &Input{
		Queue: make(chan termbox.Event),
		Ctrl:   make(chan bool, 2),
		EndKey: termbox.KeyCtrlC,
	}
}

func (input *Input) Start() {
	go poll(input)
}

func (input *Input) Stop() {
	input.Ctrl <- true
}

func poll(input *Input) {
loop:
	for {
		select {
		case <-input.Ctrl:
			break loop
		default:
			input.Queue <- termbox.PollEvent()
		}
	}
}