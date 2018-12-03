package cmd

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"path"

	_ "github.com/jicg/easypos/model"
	minit "github.com/jicg/easypos/modules/init"
	"github.com/jicg/easypos/routers"
	"github.com/jicg/easypos/routers/crm"
	"github.com/jicg/easypos/routers/crm/order"
	"github.com/jicg/easypos/routers/crm/product"
	"github.com/jicg/easypos/routers/pos"
	"github.com/jicg/easypos/routers/report"
	"github.com/jicg/easypos/routers/user"
	"gopkg.in/macaron.v1"
	"github.com/jicg/easypos/routers/crm/producttype"
	"github.com/jicg/easypos/routers/crm/muser"
)

var CmdWeb = cli.Command{
	Name:        "web",
	Usage:       "启动收银系统",
	Description: ``,
	Action:      runWeb,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 4000,
			Usage: "端口",
		},
		cli.IntFlag{
			Name:  "evn",
			Value: 0,
			Usage: "0:生产模式 ，!0：调试模式",
		},
	},
}

func runWeb(clictx *cli.Context) {
	evn := clictx.Int("evn")
	//evn = 1
	minit.Evn = evn
	if evn == 0 {
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
	}
	m := macaron.New()

	m.Map(log.New(os.Stdout, minit.DEFAULT_LOG_PREFIX, 0))
	//m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static(
		path.Join("public"),
		macaron.StaticOptions{
			Prefix:      "public",
			SkipLogging: true,
		},
	))
	m.Use(macaron.Static(
		path.Join("data/upfile"),
		macaron.StaticOptions{
			Prefix:      "data",
			SkipLogging: true,
		},
	))
	m.Use(minit.NewCache())
	m.Use(minit.NewSession())
	m.Use(minit.NewRender())
	port := clictx.Int("port")
	m.Get("", routers.Index)

	m.Get("/login", user.LoginView)
	m.Post("/login", user.Login)

	m.Group("view", func() {
		m.Get("/logout", user.Logout)
		m.Get("/changepwd", user.ChangePwd)
		m.Get("/pos", pos.Index)

		m.Group("/crm", func() {
			m.Get("/product", product.View)
			m.Get("/order", order.View)
			m.Get("/orderitem", order.ViewItem)
			m.Get("/producttype", producttype.View)
			m.Get("/user", crm.CheckAdmin,muser.View,)
		})
		m.Group("/report", func() {
			m.Get("/pos", report.PosIndex)
		})
	}, user.LoginCheck)

	//商品建档
	m.Group("/crm", func() {
		m.Group("/product", func() {
			m.Any("/query", product.Query)
			m.Any("/get/:id", product.Get)
			m.Post("/add",crm.CheckAdmin, product.Add)
			m.Post("/del/:id", crm.CheckAdmin,product.Del)
			m.Post("/edit/:id",crm.CheckAdmin, product.Edit)
			m.Any("/export", product.QueryXls)
		})
		m.Group("/producttype", func() {
			m.Any("/query", producttype.Query)
			m.Any("/get/:id", producttype.Get)
			m.Post("/add",crm.CheckAdmin, producttype.Add)
			m.Post("/del/:id", crm.CheckAdmin,producttype.Del)
			m.Post("/edit/:id",crm.CheckAdmin, producttype.Edit)
		})
		m.Group("/user", func() {
			m.Any("/query", muser.Query)
			m.Any("/get/:id", muser.Get)
			m.Post("/add", muser.Add)
			m.Post("/del/:id", muser.Del)
			m.Post("/edit/:id", muser.Edit)
		},crm.CheckAdmin)
		m.Group("/order", func() {
			m.Any("/query", order.Query)
			m.Any("/queryitem", order.QueryItem)
			m.Any("/get/:id", order.Get)
			m.Any("/getitems/:id", order.GetItems)
			m.Any("/export", order.QueryXlsOrders)
			m.Any("/exportitem", order.QueryItemXlsOrders)
		})
		m.Group("/user", func() {
			m.Post("/changepwd", user.ChangePwd2)
		})
	}, crm.LoginCheck)
	//订单
	m.Group("/pos", func() {
		m.Get("/getno", pos.GetNo)
		m.Get("/pro/:no", pos.GetPro)
		m.Get("/qpros", pos.QueryPro)

		m.Post("/create", pos.Create)
	}, crm.LoginCheck)

	m.Run(port)
}
