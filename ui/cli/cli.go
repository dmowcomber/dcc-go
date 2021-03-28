package cli

import (
	"log"

	"github.com/dmowcomber/dcc-go/rail"
	"github.com/dmowcomber/dcc-go/throttle"
	"github.com/eiannone/keyboard"
)

func New(throt *throttle.Throttle, track *rail.Track) *CLI {
	return &CLI{
		throt: throt,
		track: track,
	}
}

type CLI struct {
	throt *throttle.Throttle
	track *rail.Track
}

func (c *CLI) Run() error {
	// TODO: print usage

	err := keyboard.Open()
	if err != nil {
		return err
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		log.Printf("%q, %#v\n", char, key)

		switch char {
		case 'p':
			c.track.PowerToggle()
		case 's':
			c.throt.Stop()
		case '0':
			c.throt.ToggleFunction(0)
		case '1':
			c.throt.ToggleFunction(1)
		case '2', 'h':
			c.throt.ToggleFunction(2)
		case '3':
			c.throt.ToggleFunction(3)
		case '4':
			c.throt.ToggleFunction(4)
		case '5':
			c.throt.ToggleFunction(5)
		case '6':
			c.throt.ToggleFunction(6)
		case '7':
			c.throt.ToggleFunction(7)
		case '8', 'm':
			c.throt.ToggleFunction(8)
		case '9':
			c.throt.ToggleFunction(9)
		default:
			if key == keyboard.KeyArrowUp {
				c.throt.ThrottleUp()
			} else if key == keyboard.KeyArrowDown {
				c.throt.ThrottleDown()
			} else if key == keyboard.KeyArrowRight {
				c.throt.DirectionForward()
			} else if key == keyboard.KeyArrowLeft {
				c.throt.DirectionBackward()
			} else if key == keyboard.KeyCtrlC {
				return nil
			} else {
				// TODO: print usage
			}
		}
	}
}
