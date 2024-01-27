package app1

import (
	"fmt"
	"os"
)

func RollHWCLI() {
	args := os.Args
	fmt.Printf("Dis is de args\nos.Args: %v\n \nArgument: %v\n", args, args[1:len(args)])
}
