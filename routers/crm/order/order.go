package order

import (
	"github.com/Unknwon/com"
	"github.com/jicg/easypos/model"
	"github.com/jicg/easypos/routers/crm"
	"gopkg.in/macaron.v1"
	"strings"
)

func View(ctx *macaron.Context) {
	ctx.HTML(200, "crm/order/index")
}

func Query(ctx *macaron.Context) {
	var (
		rowCount     int
		current      int
		sort         string
		sortfelid    string
		searchPhrase string
	)
	ctx.Req.ParseForm()
	for k, v := range ctx.Req.Form {
		if strings.EqualFold(k, "rowCount") {
			rowCount = com.StrTo(v[0]).MustInt()
		}
		if strings.EqualFold(k, "current") {
			current = com.StrTo(v[0]).MustInt()
		}
		if strings.EqualFold(k, "searchPhrase") {
			searchPhrase = com.ToStr(v[0])
		}
		if strings.Contains(k, "sort") {
			sort = com.ToStr(v[0])
			a := strings.Index(k, "[") + 1
			b := strings.LastIndex(k, "]")
			sortfelid = k[a:b]
		}
	}
	session := model.Engine.NewSession()
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		session = model.Engine.Where("orderno like ?", searchPhrase).Or("c_no like ?", searchPhrase)
	}
	count, _ := session.Count(new(model.Order))

	session = model.Engine.NewSession()
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		session = model.Engine.Where("orderno like ?", searchPhrase).Or("c_no like ?", searchPhrase)
	}
	if len(sort) == 0 {
		session = session.Desc("udate")
	} else {
		if strings.EqualFold(sort, "asc") {
			session = session.Asc(sortfelid)
		} else {
			session = session.Desc(sortfelid)
		}
	}
	pros := make([]*model.Order, 0)
	session.Limit(rowCount, (current-1)*rowCount).Find(&pros)
	retData := &crm.RowDatas{
		Current:  current,
		RowCount: rowCount,
		Rows:     pros,
		Total:    count,
	}
	ctx.JSON(200, retData)
}

func Get(ctx *macaron.Context) {
	id := ctx.ParamsFloat64(":id")
	pro := new(model.Order)
	model.Engine.ID(id).Get(pro)
	ctx.JSON(200, pro)
}

func GetItems(ctx *macaron.Context) {
	id := ctx.ParamsFloat64(":id")
	pros := make([]*model.OrderItem, 0)
	model.Engine.Where("order_id = ?", id).Find(&pros)
	ctx.JSON(200, pros)
}
