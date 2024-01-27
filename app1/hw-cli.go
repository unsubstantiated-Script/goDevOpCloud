package app1

import (
	"fmt"
	"os"
)

func RollHWCLI() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Useage: ./hw-cli <argument>\n")
		os.Exit(1)
	}

	fmt.Printf("Dis is de args\nos.Args: %v\n \nArgument: %v\n", args, args[1:len(args)])
}
