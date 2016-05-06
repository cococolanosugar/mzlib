package test

import (
	"fmt"
	//"gopkg.in/mgo.v2/bson"
	"mzlib"
	"testing"
	"time"
)

func TestLogModel(t *testing.T) {
	// fmt.Println("LogModel test begin.......")
	// lm := mzlib.NewLogModel()
	// lm.Uid = "0001"
	// lm.Name = "logx"

	// lm.Put()

	// lm2 := mzlib.NewLogModel()
	// list := lm2.Find(bson.M{"name": "logx"})
	// for i, v := range list {
	// 	vv, ok := v.(*mzlib.LogModel)
	// 	fmt.Printf("list=======%d====%v====%#v========%v\n", i, ok, vv, vv.Uid)
	// }

	type LoginLog struct {
		mzlib.LogModel `bson:",inline"`
		Gold           int `json:"gold" bson:"gold"`
	}

	log := new(LoginLog)
	log.Pk = "Uid"
	log.Uid = "002"
	log.Name = "xushicai"
	log.DateTime = time.Now().Format("2006-01-02 15:04:05")
	log.Gold = 100

	mzlib.LogModelSet(log)
	//time.Sleep(3000)

	fmt.Println("LogModel test end.........")
}
