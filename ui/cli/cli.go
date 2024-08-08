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
			_, err = c.track.PowerToggle()
		case 's':
			err = c.throt.Stop()
		case '0':
			_, err = c.throt.ToggleFunction(0)
		case '1':
			_, err = c.throt.ToggleFunction(1)
		case '2', 'h':
			_, err = c.throt.ToggleFunction(2)
		case '3':
			_, err = c.throt.ToggleFunction(3)
		case '4':
			_, err = c.throt.ToggleFunction(4)
		case '5':
			_, err = c.throt.ToggleFunction(5)
		case '6':
			_, err = c.throt.ToggleFunction(6)
		case '7':
			_, err = c.throt.ToggleFunction(7)
		case '8', 'm':
			_, err = c.throt.ToggleFunction(8)
		case '9':
			_, err = c.throt.ToggleFunction(9)
		default:
			if key == keyboard.KeyArrowUp {
				err = c.throt.ThrottleUp()
			} else if key == keyboard.KeyArrowDown {
				err = c.throt.ThrottleDown()
			} else if key == keyboard.KeyArrowRight {
				err = c.throt.DirectionForward()
			} else if key == keyboard.KeyArrowLeft {
				err = c.throt.DirectionBackward()
			} else if key == keyboard.KeyCtrlC {
				return nil
			} else {
				// TODO: print usage
			}
		}
		if err != nil {
			log.Print(err.Error())
		}
	}
}
