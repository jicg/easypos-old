package cmd

import (
	"github.com/codegangsta/cli"
)

var CmdAdd = cli.Command{
	Name:        "add",
	Usage:       "自动添加模版",
	Description: ``,
	Action:      runAdd,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "action",
			Value: "",
			Usage: "动作",
		},
	},
}

func runAdd(ctx *cli.Context) {

}
