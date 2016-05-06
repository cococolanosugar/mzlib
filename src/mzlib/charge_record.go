package mzlib

import (
	"encoding/json"
	//"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ChargeRecord struct {
	Pk                  string  `json:"-" bson:"-"`
	Recharge_Id         string  `json:"recharge_id" bson:"recharge_id"`
	App_Id              string  `json:"app_id" bson:"app_id"`
	Uin                 string  `json:"uin" bson:"uin"`
	Urecharge_Id        string  `json:"urecharge_id" bson:"urecharge_id"`
	Extra               string  `json:"extra" bson:"extra"`
	Recharge_Money      float64 `json:"recharge_money" bson:"recharge_money"`
	Recharge_Gold_Count int     `json:"recharge_gold_count" bson:"recharge_gold_count"`
	Pay_Status          int     `json:"pay_status" bson:"pay_status"`
	Create_Time         string  `json:"create_time" bson:"create_time"`
	Sign                string  `json:"sign" bson:"sign"`
	Date_Time           string  `json:"date_time" bson:"date_time"`
}

func CreateChargeRecord(recharge_id string, app_id string, uin string, urecharge_id string, extra string, recharge_money float64, recharge_gold_count int, pay_status int, create_time string, sign string) *ChargeRecord {
	cr := NewChargeRecord()

	cr.Recharge_Id = recharge_id
	cr.App_Id = app_id
	cr.Uin = uin
	cr.Urecharge_Id = urecharge_id
	cr.Extra = extra
	cr.Recharge_Money = recharge_money
	cr.Recharge_Gold_Count = recharge_gold_count
	cr.Pay_Status = pay_status
	cr.Create_Time = create_time
	cr.Sign = sign

	return cr
}

func NewChargeRecord() *ChargeRecord {
	cr := new(ChargeRecord)
	cr.Pk = "Recharge_Id"
	cr.Date_Time = time.Now().Format(TIME_FORMAT)
	return cr
}

func ChargeRecordSet(rc *ChargeRecord) {
	//启动一个gorutine插入日志
	go func() {
		app.log_mongo.Insert(rc)
	}()
}

func ChargeRecordGet(recharge_id string) *ChargeRecord {
	cr := NewChargeRecord()
	//从mongo中获取
	mongo_res := app.mongo_storage.Get(cr, recharge_id)
	if mongo_res != nil {
		cr.Load(mongo_res)
		return cr
	}
	return nil
}

func (cr *ChargeRecord) Load(data interface{}) {
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := bson.Marshal(data_bson)
	bson.Unmarshal(b, cr)

}

func (cr *ChargeRecord) ToJson() interface{} {
	v, _ := json.Marshal(cr)
	return v
}

func (cr *ChargeRecord) ToBson() interface{} {
	v, _ := bson.Marshal(cr)
	return v
}
