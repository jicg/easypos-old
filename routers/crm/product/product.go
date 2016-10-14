package product

import (
	"github.com/jicg/easypos/model"
	"github.com/jicg/easypos/routers/crm"
	"gopkg.in/macaron.v1"

	"errors"
	"github.com/Unknwon/com"
	"strings"
)

func View(ctx *macaron.Context) {
	ctx.HTML(200, "crm/product/index")
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
		session = model.Engine.Where("no like ?", searchPhrase).Or("desc like ?", searchPhrase)
	}
	count, _ := session.Count(new(model.Product))

	session = model.Engine.NewSession()
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		session = model.Engine.Where("no like ?", searchPhrase).Or("desc like ?", searchPhrase)
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
	pros := make([]*model.Product, 0)
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
	pro := new(model.Product)
	model.Engine.ID(id).Get(pro)
	ctx.JSON(200, pro)
}

func Add(ctx *macaron.Context) {
	no := ctx.Query("no")
	desc := ctx.QueryTrim("desc")
	price := ctx.QueryFloat64("price")
	saleprice := ctx.QueryFloat64("saleprice")
	qtycan := ctx.QueryFloat64("qtycan")
	err := model.AddProduct(&model.Product{No: no, Desc: desc, Price: price, Saleprice: saleprice, Qtycan: qtycan})
	jsonRet := &crm.RetJson{Code: 0, Msg: "添加成功！"}
	if err != nil {
		jsonRet.Code = -1
		jsonRet.Msg = err.Error()
	}
	ctx.JSON(200, jsonRet)
}

func Del(ctx *macaron.Context) {
	id := ctx.ParamsInt64(":id")
	_, err := model.Engine.Where(" id = ? ", id).Delete(new(model.Product))
	jsonRet := &crm.RetJson{Code: 0, Msg: "删除成功！"}
	if err != nil {
		jsonRet.Code = -1
		jsonRet.Msg = err.Error()
	}
	ctx.JSON(200, jsonRet)
}

func Edit(ctx *macaron.Context) {
	id := ctx.ParamsInt64(":id")
	pro := new(model.Product)
	model.Engine.ID(id).Get(pro)

	var err error
	if pro != nil {
		desc := ctx.QueryTrim("desc")
		price := ctx.QueryFloat64("price")
		saleprice := ctx.QueryFloat64("saleprice")
		qtycan := ctx.QueryFloat64("qtycan")

		pro.Desc = desc
		pro.Price = price
		pro.Saleprice = saleprice
		pro.Qtycan = qtycan

		_, err = model.Engine.ID(id).Update(pro)
	} else {
		err = errors.New("商品不存在")
	}

	jsonRet := &crm.RetJson{Code: 0, Msg: "修改成功"}
	if err != nil {
		jsonRet.Code = -1
		jsonRet.Msg = err.Error()
	}
	ctx.JSON(200, jsonRet)
}
