package main

import (
	"fmt"
	"os"

	"github.com/pedrosantan4/egocli/cmd"
)

func main() {
	// Ensure we're in the right directory
	if _, err := os.Stat("go.mod"); err != nil {
		fmt.Println("⚠️ WARNING: Not running in project root")
	}
	cmd.Execute()
}
