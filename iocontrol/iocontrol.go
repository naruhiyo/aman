package iocontrol

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nsf/termbox-go"
)

func DeleteInput(inputs *string) {
	var space = ""
	for i := 0; i < len(*inputs); i++ {
		space += " "
	}
	fmt.Printf("\r%s", space)
	if (0 < len(*inputs)) {
		*inputs = (*inputs)[:len(*inputs)-1]
	}
}

func ReceiveKeys(inputs *string) int {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
			case termbox.KeyEsc:
				return -1
			case termbox.KeySpace:
				*inputs += " "
			case termbox.KeyBackspace:
				DeleteInput(inputs)
			case termbox.KeyBackspace2:
				DeleteInput(inputs)
			default:
				*inputs += string(ev.Ch)
		}
	default:
	}
	return 0
}

func RenderQuery(inputs *string) {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	fmt.Printf("\r%s\n", *inputs)
}

func RenderResult(result []string) {
	for i := 0; i < len(result); i++ {
		fmt.Printf("\r%s\n", result[i])
	}
}
