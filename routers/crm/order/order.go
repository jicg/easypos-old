package order

import (
	"github.com/Unknwon/com"
	"github.com/jicg/easypos/model"
	"github.com/jicg/easypos/routers/crm"
	"gopkg.in/macaron.v1"
	"strings"
	"github.com/xuri/excelize"
	"fmt"
	"github.com/guregu/null"
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
		session = model.Engine.Where("orderno like ?", searchPhrase).Or("customno like ?", searchPhrase)
	}
	count, _ := session.Count(new(model.Order))

	session = model.Engine.NewSession()
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		session = model.Engine.Where("orderno like ?", searchPhrase).Or("customno like ?", searchPhrase)
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
//
func QueryXlsOrders(ctx *macaron.Context) {
	var (
		searchPhrase string
	)
	ctx.Req.ParseForm()
	for k, v := range ctx.Req.Form {
		if strings.EqualFold(k, "searchPhrase") {
			searchPhrase = com.ToStr(v[0])
		}
	}

	//session := model.Engine.NewSession()
	var (
		ret []map[string][]byte
		err error
	)

	//pros := make([]*model.Order, 0)
	//session.Find(&pros)

	sql := `
	select t.id, tt.orderno,tt.customno,tt.user_name,t.product_desc,
  		t.qty,t.trueprice price ,tt.desc,t.cdate
	from ` + "`order`" + ` tt,order_item  t
	where tt.id = t.order_id %s
	order by t.cdate desc
	`
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		ssql := " and (orderno like ? or customno like ?)"
		sql = fmt.Sprintf(sql, ssql)
		ret, err = model.Engine.Query(sql, searchPhrase, searchPhrase)
	} else {
		ret, err = model.Engine.Query(fmt.Sprintf(sql, ""))
	}

	if err != nil {
		ctx.JSON(200, &crm.RetJson{Code: -1, Msg: "错误：" + err.Error(), Data: nil})
		return
	}
	xlsx := excelize.NewFile()
	if style, err := xlsx.NewStyle(`{"font":{"bold":true,family":"Berlin Sans FB Demi"}}`); err == nil {
		xlsx.SetCellStyle("Sheet1", "A1", "I1", style)
	}

	xlsx.SetCellValue("Sheet1", "A1", "Id")
	xlsx.SetCellValue("Sheet1", "B1", "订单编号")
	xlsx.SetCellValue("Sheet1", "C1", "手工单号")
	xlsx.SetCellValue("Sheet1", "D1", "员工")
	xlsx.SetCellValue("Sheet1", "E1", "商品")
	xlsx.SetCellValue("Sheet1", "F1", "数量")
	xlsx.SetCellValue("Sheet1", "G1", "单价")
	xlsx.SetCellValue("Sheet1", "H1", "备注")
	xlsx.SetCellValue("Sheet1", "I1", "创建时间")
	for index, v := range ret {
		indexStr := index + 2
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "A", indexStr), v["id"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "B", indexStr), v["orderno"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "C", indexStr), v["customno"])

		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "D", indexStr), v["user_name"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "E", indexStr), v["product_desc"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "F", indexStr), v["qty"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "G", indexStr), v["price"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "H", indexStr), v["desc"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "I", indexStr), v["cdate"])
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
		"attachment;filename=orders.xlsx",
	}
	ctx.Header()["Content-Transfer-Encoding"]=[]string{
		"binary",
	}
	xlsx.Write(ctx.Resp)
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


func ViewItem(ctx *macaron.Context)  {
	ctx.HTML(200, "crm/order/item")
}
func QueryItem(ctx *macaron.Context) {
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
	type RetS struct {
		Id int `json:"id"`
		Orderno null.String `json:"orderno"`
		Customno  null.String `json:"customno"`
		UserName  null.String `db:"user_name" json:"user_name"`
		ProductDesc  null.String `db:"product_desc"  json:"product_desc"`
		Qty  null.Float `json:"qty"`
		Price  null.Float `json:"price"`
		Desc  null.String `json:"desc"`
		Udate  model.Time `json:"udate"`
	}
	var (
		data []RetS
		err error
	)
	sqlcount := `
	select count(1)
	from ` + "`order`" + ` tt,order_item  t
	where tt.id = t.order_id %s
	`
	total :=int64(0)
	var param []interface{}
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		ssql := " and (orderno like ? or customno like ? or product_desc like ? or product_no like ? )"
		sqlcount  =fmt.Sprintf(sqlcount, ssql)
		param=[]interface{}{searchPhrase,searchPhrase,searchPhrase,searchPhrase}
	} else {
		sqlcount=fmt.Sprintf(sqlcount, "")
		param=[]interface{}{}
	}
	err =model.SqlxDB.Get(&total,sqlcount,param...)
	if err!=nil{
		ctx.JSON(200, &crm.RetJson{Code: -1, Msg: "错误：" + err.Error(), Data: nil})
		return
	}



	sql := `
select id ,orderno,customno,user_name,product_desc,qty,price,desc,udate
from (
	select t.id, tt.orderno,tt.customno,tt.user_name,t.product_desc,
  		t.qty,t.trueprice price ,tt.desc,t.udate
	from ` + "`order`" + ` tt,order_item  t
	where tt.id = t.order_id %s)
	`
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		sql = fmt.Sprintf(sql, " and (orderno like ? or customno like ? or product_desc like ? or product_no like ? ) ")
		param = []interface{}{ searchPhrase, searchPhrase,searchPhrase,searchPhrase}
	} else {
		sql =fmt.Sprintf(sql, "" )
		param = []interface{}{}
	}
	sqlorder:=" order by udate desc "
	if len(sortfelid)!=0{
		sqlorder=fmt.Sprintf(" order by %s %s " ,sortfelid,sort)
	}
	sql+=sqlorder
	sql+=fmt.Sprintf(" limit %d, %d",(current-1)*rowCount,rowCount)
	err =model.SqlxDB.Select(&data,sql, param ...)
	if err != nil {
		ctx.JSON(200, &crm.RetJson{Code: -1, Msg: "错误：" + err.Error(), Data: nil})
		return
	}

	if err != nil {
		ctx.JSON(200, &crm.RetJson{Code: -1, Msg: "错误：" + err.Error(), Data: nil})
		return
	}
	retData := &crm.RowDatas{
		Current:  current,
		RowCount: rowCount,
		Rows:     data,
		Total:    total,
	}
	ctx.JSON(200, retData)
}

func QueryItemXlsOrders(ctx *macaron.Context) {
	var (
		searchPhrase string
	)
	ctx.Req.ParseForm()
	for k, v := range ctx.Req.Form {
		if strings.EqualFold(k, "searchPhrase") {
			searchPhrase = com.ToStr(v[0])
		}
	}

	//session := model.Engine.NewSession()
	var (
		ret []map[string][]byte
		err error
	)
	sql := `
select id ,orderno,customno,user_name,product_desc,qty,price,desc,udate
from (
	select t.id, tt.orderno,tt.customno,tt.user_name,t.product_desc,
  		t.qty,t.trueprice price ,tt.desc,t.udate
	from ` + "`order`" + ` tt,order_item  t
	where tt.id = t.order_id %s)
	`
	param := []interface{}{}
	if len(searchPhrase) != 0 {
		searchPhrase = "%" + searchPhrase + "%"
		sql = fmt.Sprintf(sql, " and (orderno like ? or customno like ? or product_desc like ? or product_no like ? ) ")
		param = []interface{}{ searchPhrase, searchPhrase,searchPhrase,searchPhrase}
	} else {
		sql =fmt.Sprintf(sql, "" )
		param = []interface{}{}
	}
	ret, err = model.Engine.Query(sql,param...)

	if err != nil {
		ctx.JSON(200, &crm.RetJson{Code: -1, Msg: "错误：" + err.Error(), Data: nil})
		return
	}
	xlsx := excelize.NewFile()
	if style, err := xlsx.NewStyle(`{"font":{"bold":true,family":"Berlin Sans FB Demi"}}`); err == nil {
		xlsx.SetCellStyle("Sheet1", "A1", "I1", style)
	}

	xlsx.SetCellValue("Sheet1", "A1", "Id")
	xlsx.SetCellValue("Sheet1", "B1", "订单编号")
	xlsx.SetCellValue("Sheet1", "C1", "手工单号")
	xlsx.SetCellValue("Sheet1", "D1", "员工")
	xlsx.SetCellValue("Sheet1", "E1", "商品")
	xlsx.SetCellValue("Sheet1", "F1", "数量")
	xlsx.SetCellValue("Sheet1", "G1", "单价")
	xlsx.SetCellValue("Sheet1", "H1", "备注")
	xlsx.SetCellValue("Sheet1", "I1", "修改时间")
	for index, v := range ret {
		indexStr := index + 2
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "A", indexStr), v["id"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "B", indexStr), v["orderno"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "C", indexStr), v["customno"])

		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "D", indexStr), v["user_name"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "E", indexStr), v["product_desc"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "F", indexStr), v["qty"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "G", indexStr), v["price"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "H", indexStr), v["desc"])
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "I", indexStr), v["udate"])
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
		"attachment;filename=orders.xlsx",
	}
	ctx.Header()["Content-Transfer-Encoding"]=[]string{
		"binary",
	}
	xlsx.Write(ctx.Resp)
}

