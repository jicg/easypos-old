package pos

import (
	"encoding/json"
	"fmt"
	"github.com/jicg/easypos/model"
	"gopkg.in/macaron.v1"
	"math/rand"
	"strconv"
	"time"
)

//POS201405060921666AAA

// 获取 零售单号 /
// 查询 条码 /no

//新增订单 json

func Index(ctx *macaron.Context) {
	ctx.HTML(200, "pos/index")
}

func GetNo(ctx *macaron.Context) {
	bb := "P" + time.Now().Format("20060102150405")
	for i := 0; i < 3; i++ {
		a := rand.Intn(3)
		if a == 0 {
			bb = bb + strconv.Itoa(rand.Intn(10))
		} else {
			b := int('a')
			if a == 1 {
				b = int('a') + rand.Intn(26)
			} else {
				b = int('A') + rand.Intn(26)
			}
			bb = bb + string(rune(b))
		}
	}
	data := &RetData{Code: 0, Msg: "ok", Data: bb}
	ctx.JSON(200, data) //??
}

func GetPro(ctx *macaron.Context) {
	no := ctx.Params(":no")
	pro := new(model.Product)
	_, err := model.Engine.Where("no like ?", no+"%").Limit(1).Get(pro)

	data := &RetData{Code: 0, Msg: "ok", Data: pro}
	if err != nil {
		data.Code = -1
		data.Msg = err.Error()
	}
	if pro.Id == 0 {
		data.Code = -1
		data.Msg = "条码不存在"
	}
	ctx.JSON(200, data)
}

func QueryPro(ctx *macaron.Context) {
	no := ctx.Query("query")

	type Se struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	pros := make([]*model.Product, 0)
	model.Engine.Where(
		"no like ? ",
		"%"+no+"%",
	).Or(
		"desc like ? ",
		"%"+no+"%",
	).Limit(10).Find(&pros)
	ss := make([]*Se, len(pros))
	for index, p := range pros {
		ss[index] = &Se{p.Desc, p.No}
	}
	ctx.JSON(200, map[string]interface{}{
		"suggestions": ss,
	})
}

func Create(ctx *macaron.Context) {
	querystr := ctx.QueryTrim("data")
	var porder model.Order
	ret := &RetData{Code: 0, Msg: "ok"}
	json.Unmarshal([]byte(querystr), &porder)

	if &porder.Items == nil || len(porder.Items) == 0 {
		ret.Code = -1
		ret.Msg = "明细不能为空！"
		ctx.JSON(200, ret)
		return
	}
	tnow := time.Now()
	porder.Cdate = model.Time(tnow)
	porder.Udate = model.Time(tnow)

	user := ctx.Data["user"]
	u := user.(*model.User)
	porder.UserId = int(u.Id)
	porder.UserName = u.Name

	session := model.Engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		ctx.JSON(200, ret)
		return
	}
	if _, err := session.InsertOne(&porder); err != nil {
		session.Rollback()
		ret.Code = -1
		ret.Msg = err.Error()
		ctx.JSON(200, ret)
		return
	}

	for _, item := range porder.Items {
		item.Cdate = model.Time(tnow)
		item.Udate = model.Time(tnow)
		item.OrderId = porder.Id
	}

	if _, err := session.Insert(&porder.Items); err != nil {
		session.Rollback()
		ret.Code = -1
		ret.Msg = err.Error()
		ctx.JSON(200, ret)
		return
	}

	//减少库存
	for _, item := range porder.Items {
		var pro = new(model.Product)
		if session.Where("no = ?", item.ProductNo).Get(pro); pro != nil {
			pro.Qtycan = pro.Qtycan - item.Qty
			//if pro.Qtycan < 0 {
			//	session.Rollback()
			//	ret.Code = -1
			//	ret.Msg = fmt.Sprintf("条码[%s]库存不够,当前库存为: %.2f ", item.ProductDesc, pro.Qtycan+item.Qty)
			//	ctx.JSON(200, ret)
			//	return
			//}
			if _, err := session.Id(pro.Id).Update(pro); err != nil {
				session.Rollback()
				ret.Code = -1
				ret.Msg = fmt.Sprintf("条码[%s]更新库存失败,原因:%s", item.ProductDesc, err.Error())
				ctx.JSON(200, ret)
				return
			}
		} else {
			session.Rollback()
			ret.Code = -1
			ret.Msg = fmt.Sprintf("条码[%s]不存在", item.ProductDesc)
			ctx.JSON(200, ret)
			return
		}
	}
	session.Commit()
	ctx.JSON(200, ret)
}
