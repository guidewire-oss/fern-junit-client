package main

import (
	"fmt"
	"os"

	"github.com/guidewire-oss/fern-junit-client/cmd"
)

func main() {
	banner, err := os.ReadFile("./static/banner.txt")
	if err == nil {
		fmt.Println(string(banner))
	}
	cmd.Execute()
}
