package report

import (
	"github.com/jicg/easypos/model"
	"gopkg.in/macaron.v1"
	"log"
	"strconv"
)

type PosSale struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

func PosIndex(ctx *macaron.Context) {
	rets, err := model.Engine.Query(`select no,totamt,proamt from 
		(select a.product_no as no,sum(a.trueprice*a.qty) as totamt ,sum((a.trueprice-b.price)*a.qty) as proamt 
		from order_item a ,product b 
		where strftime("%Y-%m",a.cdate) = strftime("%Y-%m",'now')
		and a.product_no = b.no
		group by a.product_no ) order by proamt desc limit 8 
	`)
	if err != nil {
		log.Println("错误：%v\n", err)
	}
	nos := make([]string, 0)
	totamts := make([]float64, 0)
	proamts := make([]float64, 0)
	for i := 0; i < len(rets); i++ {
		ret := rets[i]
		nos = append(nos, string(ret["no"]))
		v1, err := strconv.ParseFloat(string(ret["totamt"]), 64)
		if err != nil {
			v1 = float64(0)
		}
		totamts = append(totamts, v1)

		v2, err := strconv.ParseFloat(string(ret["proamt"]), 64)
		if err != nil {
			v2 = float64(0)
		}
		proamts = append(proamts, v2)
	}
	ctx.Data["categories"] = nos
	ctx.Data["datas"] = [...]*PosSale{{Name: "总金额", Data: totamts}, {Name: "收益金额", Data: proamts}}
	ctx.HTML(200, "report/index")
}
