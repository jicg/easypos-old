package routers

import (
	"gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) {
	ctx.Redirect("/view/pos")
}
