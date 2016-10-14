package model

import (
	"time"
)

type OrderItem struct {
	Id          int64   `json:"id"           xorm:"pk autoincr"`
	OrderId     int64   `json:"order_id"     xorm:"not null index"`
	ProductId   int64   `json:"product_id"`
	ProductNo   string  `json:"product_no"   xorm:"default '' VARCHAR(80) index"`
	ProductDesc string  `json:"product_desc" xorm:"default '' TEXT"`
	Saleprice   float64 `json:"saleprice"    xorm:" default 0 Float"`
	Trueprice   float64 `json:"trueprice"    xorm:"default 0 Float"`
	Qty         float64 `json:"qty"          xorm:"default 0"`
	Cdate       Time    `json:"cdate"        xorm:"created"`
	Udate       Time    `json:"udate"        xorm:"updated"`
}

type Order struct {
	Id       int64        `json:"id"       xorm:"not null pk autoincr"`
	Orderno  string       `json:"orderno"  xorm:"not null unique"`
	Customno string       `json:"customno" xorm:"index default ''"`
	Desc     string       `json:"desc"   xorm:"TEXT"`
	Totamt   float64      `json:"totamt"   xorm:"default 0.00 Float"`
	Trueamt  float64      `json:"trueamt"  xorm:"default 0.00 Float"`
	Payamt   float64      `json:"payamt"   xorm:"Float"`
	Retamt   float64      `json:"retamt"   xorm:"default 0.00 Float"`
	Cdate    Time         `json:"cdate"    xorm:"created"`
	Udate    Time         `json:"udate"    xorm:"updated index"`
	Items    []*OrderItem `json:"items"    xorm:"-"`
}

type Product struct {
	Id        int64   `json:"id"         xorm:"not null pk autoincr"`
	No        string  `json:"no"         xorm:"not null VARCHAR(80) unique"`
	Desc      string  `json:"desc"     xorm:"not null TEXT index"`
	Picurl    string  `json:"picurl"     xorm:"TEXT"`
	Price     float64 `json:"price"      xorm:"default 0 Float"`
	Saleprice float64 `json:"saleprice"  xorm:"default 0 Float"`
	Qtycan    float64 `json:"qtycan"     xorm:"default 0"`
	Cdate     Time    `json:"cdate"      xorm:"created"`
	Udate     Time    `json:"udate"      xorm:"updated"`
}

type User struct {
	Id       int64  `xorm:"not null pk autoincr"`
	Name     string `xorm:"not null VARCHAR(80) unique"`
	Nickname string
	Pwd      string `xorm:"not null VARCHAR(80)"`
	Rands    string `xorm:"VARCHAR(10)"`
	Salt     string `xorm:"VARCHAR(10)"`
	IsAdmin  bool
	Cdate    Time `xorm:"created"`
	Udate    Time `xorm:"updated"`
}

type Time time.Time

const (
	timeFormart = "20060102 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) GobEncode() ([]byte, error) {
	return time.Time(t).MarshalBinary()
}

func (t *Time) GobDecode(data []byte) error {
	var lt time.Time
	if err := lt.UnmarshalBinary(data); err != nil {
		return err
	} else {
		*t = Time(lt)
		return nil
	}
}
