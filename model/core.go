package model

import (
	"encoding/gob"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	mini "github.com/jicg/easypos/modules/init"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var Engine *xorm.Engine

func init() {

	log.Println("＊＊＊＊＊＊＊＊＊数据库初始化＊＊＊＊＊＊＊＊＊＊＊")
	//得到配置信息
	var err error
	Engine, err = xorm.NewEngine("sqlite3", "./data/easypos.db")
	if err != nil {
		println(err.Error())
		return
	}
	f, err := os.Open(mini.XORM_LOG_PATH)
	if err != nil {
		f, err = os.Create(mini.XORM_LOG_PATH)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	Engine.SetLogger(xorm.NewSimpleLogger(f))
	if mini.Evn != 0 {
		Engine.ShowSQL(true)
	} else {
		Engine.Logger().SetLevel(core.LOG_INFO)
	}

	//缓存
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	Engine.SetDefaultCacher(cacher)
	log.Println("＊＊＊＊＊＊＊＊＊同步表结构＊＊＊＊＊＊＊＊＊＊＊")

	err = Engine.Sync2(new(Order), new(OrderItem), new (Producttype),new(Product), new(User))
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Println("＊＊＊＊＊＊＊＊＊数据库加载成功＊＊＊＊＊＊＊＊＊＊＊")
	gob.Register(&Time{})
	gob.Register(&User{})
	if cnt, err := GetUserCount(); err == nil && cnt == 0 {
		Engine.InsertOne(&User{Name: "admin", Nickname: "管理员", Pwd: "admin", IsAdmin: true})
	}
}
