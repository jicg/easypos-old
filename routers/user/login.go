package user

import (
	"fmt"
	"github.com/go-macaron/session"
	"github.com/jicg/easypos/model"
	"github.com/jicg/easypos/routers/crm"
	"gopkg.in/macaron.v1"
	"strings"
)

func LoginView(ctx *macaron.Context) {
	ctx.HTML(200, "login")
}

func ChangePwd(ctx *macaron.Context) {
	ctx.HTML(200, "user/changepwd")
}
func ChangePwd2(ctx *macaron.Context, sess session.Store) {
	oldpass := ctx.QueryTrim("oldpass")
	newpass := ctx.QueryTrim("newpass")
	if len(newpass) != 0 && len(oldpass) != 0 {
		user := (ctx.Data["user"]).(*model.User)
		if strings.EqualFold(oldpass, user.Pwd) {
			user.Pwd = newpass
			if _, err := model.Engine.Id(user.Id).Update(user); err != nil {
				ctx.JSON(200, &crm.RetJson{Code: -1, Msg: fmt.Sprintf("密码更新失败，失败原因:[%s]！", err.Error())})
			} else {
				sess.Set("USER", user)
				ctx.JSON(200, &crm.RetJson{Code: 0, Msg: "ok"})
			}
		} else {
			ctx.JSON(200, &crm.RetJson{Code: -1, Msg: "原密码不对！"})
		}

	} else {
		ctx.JSON(200, &crm.RetJson{Code: -1, Msg: "旧密码和新密码都不能为空！"})
	}
}

func Login(ctx *macaron.Context, sess session.Store) {
	user := sess.Get("USER")
	if user == nil {
		name := ctx.QueryTrim("username")
		pwd := ctx.QueryTrim("password")
		if len(name) != 0 && len(pwd) != 0 {
			user := model.GetUserByNameWithPwd(name, pwd)
			if user.Id == 0 {
				ctx.Data["error"] = "用户名或密码错误"
			} else {
				sess.Set("USER", user)
				ctx.Redirect("/view/pos")
				return
			}
		}
	} else {
		ctx.Data["user"] = user
	}
	ctx.Redirect("/view/pos")
}

func Logout(ctx *macaron.Context, sess session.Store) {
	sess.Destory(ctx)
	ctx.Redirect("/login")
}

func LoginCheck(ctx *macaron.Context, sess session.Store) {
	user := sess.Get("USER")
	if user == nil {
		ctx.Redirect("/login")
	} else {
		ctx.Data["user"] = user
	}
}
