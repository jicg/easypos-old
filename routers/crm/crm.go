package crm

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"github.com/jicg/easypos/model"
)

type RetJson struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type RowDatas struct {
	Current  int         `json:"current"`
	RowCount int         `json:"rowCount"`
	Rows     interface{} `json:"rows"`
	Total    int64       `json:"total"`
}

func LoginCheck(ctx *macaron.Context, sess session.Store) {
	user := sess.Get("USER")
	if user == nil {
		ctx.JSON(200, &RetJson{Code: -1, Msg: "用户没登录！"})
		return
	} else {
		ctx.Data["user"] = user
	}
}

func CheckAdmin(ctx *macaron.Context) {
	user := ctx.Data["user"]
	if user == nil {
		ctx.JSON(200, &RetJson{Code: -1, Msg: "用户没登录！"})
		return
	}
	u := user.(*model.User)
	if !u.IsAdmin {
		ctx.JSON(200, &RetJson{Code: -1, Msg: "权限不足！"})
		return
	}
}
