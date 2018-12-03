package product

import (
	"github.com/jicg/easypos/model"
	"github.com/jicg/easypos/routers/crm"
	"gopkg.in/macaron.v1"

	"errors"
	"github.com/Unknwon/com"
	"strings"
	"github.com/jicg/easypos/modules/log"
	"fmt"
	"github.com/xuri/excelize"
)

func View(ctx *macaron.Context) {
	ds := make([]*model.Producttype, 0)
	if err := model.Engine.Find(&ds); err != nil {
		log.Error(0, "错误:%s", err.Error())
	}
	ctx.Data["ptypes"] = ds
	ctx.HTML(200, "crm/product/index")
}

func Query(ctx *macaron.Context) {
	var (
		rowCount     int
		current      int
		sort         string
		sortfelid    string
		searchPhrase string
		producttype  int
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
		if strings.EqualFold(k, "producttype") {
			producttype = com.StrTo(v[0]).MustInt()
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
		session = model.Engine.Where("(no like ? or desc like ?)", searchPhrase,searchPhrase)
	}
	if producttype != 0 {
		session.And("producttype_id = ?", producttype)
	}
	count, _ := session.Count(new(model.Product))

	session = model.Engine.NewSession()
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		session = model.Engine.Where("(no like ? or desc like ?)", searchPhrase,searchPhrase)
	}
	if producttype != 0 {
		session.And("producttype_id = ?", producttype)
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


func QueryXls(ctx *macaron.Context) {
	var (
		searchPhrase string
		producttype  int
	)
	ctx.Req.ParseForm()
	for k, v := range ctx.Req.Form {

		if strings.EqualFold(k, "searchPhrase") {
			searchPhrase = com.ToStr(v[0])
		}
		if strings.EqualFold(k, "producttype") {
			producttype = com.StrTo(v[0]).MustInt()
		}
	}
	session := model.Engine.NewSession()
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		session = model.Engine.Where("(no like ? or desc like ?)", searchPhrase,searchPhrase)
	}
	if producttype != 0 {
		session.And("producttype_id = ?", producttype)
	}

	pros := make([]*model.Product, 0)
	err:=session.Find(&pros)
	if err!=nil{
		ctx.JSON(200,&crm.RetJson{Code:-1,Msg:"错误："+err.Error()})
		return
	}
	xlsx := excelize.NewFile()
	if style, err := xlsx.NewStyle(`{"font":{"bold":true,family":"Berlin Sans FB Demi"}}`); err == nil {
		xlsx.SetCellStyle("Sheet1", "A1", "I1", style)
	}
	xlsx.SetCellValue("Sheet1", "A1", "Id")
	xlsx.SetCellValue("Sheet1", "B1", "商品名称")
	xlsx.SetCellValue("Sheet1", "C1", "商品类别")
	xlsx.SetCellValue("Sheet1", "D1", "进价")
	xlsx.SetCellValue("Sheet1", "E1", "售价")
	xlsx.SetCellValue("Sheet1", "F1", "库存")
	xlsx.SetCellValue("Sheet1", "G1", "单位")
	xlsx.SetCellValue("Sheet1", "H1", "备注")
	xlsx.SetCellValue("Sheet1", "I1", "创建时间")
	for index, v := range pros {
		indexStr := index + 2
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "A", indexStr), v.Id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "B", indexStr), v.Desc)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "C", indexStr), v.ProducttypeName)

		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "D", indexStr), v.Price)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "E", indexStr), v.Saleprice)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "F", indexStr), v.Qtycan)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "G", indexStr), v.Unit)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "H", indexStr), v.Remark)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "I", indexStr), v.Cdate.ToDateString())
	}
	xlsx.SetActiveSheet(1)
	ctx.Header()["Expires"]=[]string{"0"}
	ctx.Header()["Cache-Control"] = []string{
		"must-revalidate, post-check=0, pre-check=0",
	}
	ctx.Header()["Content-Type"] = []string{
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/download",
		"application/octet-stream",
		"application/force-download",
	}
	ctx.Header()["Content-Disposition"] = []string{
		"attachment;filename=product.xlsx",
	}
	ctx.Header()["Content-Transfer-Encoding"]=[]string{
		"binary",
	}
	xlsx.Write(ctx.Resp)
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
	remark := ctx.QueryTrim("remark")
	unit := ctx.QueryTrim("unit")
	producttype_id := ctx.QueryInt("producttype_id")
	pro := new(model.Producttype)

	jsonRet := &crm.RetJson{Code: 0, Msg: "添加成功！"}
	if _, err := model.Engine.Id(producttype_id).Get(pro); err != nil {
		if err != nil {
			jsonRet.Code = -1
			jsonRet.Msg = err.Error()
		}
		ctx.JSON(200, jsonRet)
		return
	}

	if err := model.AddProduct(
		&model.Product{
			No: no,
			Desc: desc,
			Price: price,
			Saleprice: saleprice,
			Qtycan: qtycan,
			ProducttypeName:pro.Name,
			ProducttypeId:producttype_id,
			Unit:unit,
			Remark: remark,
		}); err != nil {
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
	_,err:=model.Engine.ID(id).Get(pro)

	if err == nil {
		desc := ctx.QueryTrim("desc")
		no := ctx.QueryTrim("no")
		producttype_id := ctx.QueryInt("producttype_id")
		price := ctx.QueryFloat64("price")
		saleprice := ctx.QueryFloat64("saleprice")
		qtycan := ctx.QueryFloat64("qtycan")
		remark := ctx.QueryTrim("remark")
		unit := ctx.QueryTrim("unit")

		pt := new(model.Producttype)
		_,err= model.Engine.ID(producttype_id).Get(pt)
		if err==nil {
			pro.Desc = desc
			pro.Price = price
			pro.Saleprice = saleprice
			pro.Qtycan = qtycan
			pro.ProducttypeId = int(pt.Id)
			pro.ProducttypeName = pt.Name
			pro.No = no
			pro.Unit = unit
			pro.Remark = remark
		}else{
			err = errors.New("商品类别不存在 "+err.Error())
		}
		_, err = model.Engine.ID(id).MustCols("remark,qty,unit,price,saleprice,qtycan").Update(pro)
	} else {
		err = errors.New("商品不存在 "+err.Error())
	}

	jsonRet := &crm.RetJson{Code: 0, Msg: "修改成功"}
	if err != nil {
		jsonRet.Code = -1
		jsonRet.Msg = err.Error()
	}
	ctx.JSON(200, jsonRet)
}
