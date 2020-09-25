package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aman/modules"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()
	fmt.Println(text)
	modules.Hello()
}
