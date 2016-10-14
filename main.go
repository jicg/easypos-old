package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jicg/easypos/cmd"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmd.CmdWeb,
		cmd.CmdAdd,
	}
	if len(os.Args) <= 1 {
		os.Args = []string{"./easypos", "web"}
	}
	app.Run(os.Args)
}
