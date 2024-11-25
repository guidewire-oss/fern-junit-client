package main

import (
	"fmt"
	"strings"

	"github.com/guidewire-oss/fern-junit-client/cmd"
)

const banner = `
 _____                     _ _   _       _ _      ____ _ _            _
|  ___|__ _ __ _ __       | | | | |_ __ (_) |_   / ___| (_) ___ _ __ | |_
| |_ / _ \ '__| '_ \   _  | | | | | '_ \| | __| | |   | | |/ _ \ '_ \| __|
|  _|  __/ |  | | | | | |_| | |_| | | | | | |_  | |___| | |  __/ | | | |_
|_|  \___|_|  |_| |_|  \___/ \___/|_| |_|_|\__|  \____|_|_|\___|_| |_|\__|

`

func main() {
	fmt.Print(strings.TrimLeft(banner, "\n"))
	cmd.Execute()
}
