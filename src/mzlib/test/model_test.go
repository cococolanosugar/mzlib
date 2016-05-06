package test

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"mzlib"
	"testing"
)

func TestBaseModel(t *testing.T) {

	fmt.Println("BaseModel test begin.......")

	info := simplejson.New()
	info.Set("name", "xushicai")

	bm := mzlib.BaseModel{
		Pk:  "Uid",
		Uid: "15618973660",
		Info: map[string]interface{}{
			"ppp": "7777",
		},
	}
	bm.Put()
	bm.DoPut()

	bm2 := new(mzlib.BaseModel)
	bm2.Pk = "Uid"
	bm2.Get("15618973660")
	v2 := bm2.ToJson()
	fmt.Printf("TestModel v2=======%s\n", v2)

	bm3 := new(mzlib.BaseModel)
	bm3.Get("15618973660")
	fmt.Printf("bm3=======%s\n", bm3.ToJson())

	bm4 := new(mzlib.BaseModel)
	bm4.Load([]byte(`{"uid":"15618973660","info":{"ppp":"7777"}}`))
	fmt.Printf("bm4=======%s\n", bm4.ToJson())

	bm2.Delete()

	fmt.Println("BaseModel test end.........")
}

func TestUserModel(t *testing.T) {
	fmt.Println("UserModel test begin.......")

	um0 := new(mzlib.UserModel)
	um0.Pk = "Uid"
	um0.Uid = "15618973660"
	um0.Info = map[string]interface{}{
		"ggg": "uuuu",
	}
	um0.Name = "xushicai"
	um0.Put()

	um := new(mzlib.UserModel)
	um.Pk = "Uid"

	um2 := um.Get("15618973660").(*mzlib.UserModel)
	fmt.Printf("UserModel um2 ToJson=======%s\n", um2.ToJson())
	fmt.Printf("UserModel um2.Uid=======%s\n", um2.Uid)
	fmt.Printf("UserModel um2 ToBson=======%s\n", um2.ToBson())
	um2.Name = "ppp"
	um2.DoPut()
	um2.Delete()

	type UserCards struct {
		mzlib.UserModel `bson:",inline"`
		Cards           []string `json:"cards" bson:"cards"`
	}

	uc := new(UserCards)
	uc.Pk = "Uid"
	uc.Uid = "166666666"
	uc.Name = "xushicai"
	uc.Info = map[string]interface{}{
		"yyyy": "xxxx",
	}
	uc.Cards = []string{"card_1", "card_2"}
	uc.Put()
	uc.DoPut()

	uo := new(UserObject)
	uo.Uid = "uoid_111"
	uo.Name = "uo"

	po := new(PlayerObject)
	po.Uid = "poid_111"
	po.Name = "po"
	po.Heros = []string{"hero_1", "hero_2"}

	uo_json_1, _ := json.Marshal(uo)
	uo_json_2 := uo.ToJson()
	po_json_1, _ := json.Marshal(po)
	po_json_2 := po.ToJson()

	fmt.Printf("uo_json_1=========%s\n", uo_json_1)
	fmt.Printf("uo_json_2=========%s\n", uo_json_2)
	fmt.Printf("po_json_1=========%s\n", po_json_1)
	fmt.Printf("po_json_2=========%s\n", po_json_2)

	var ui User = po
	fmt.Printf("ui=========%s\n", ui.ToJson())

	fmt.Println("UserModel test end.........")
}

type User interface {
	ToJson() []byte
}

type UserObject struct {
	Uid  string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
}

func (uo *UserObject) ToJson() []byte {
	js, _ := json.Marshal(uo)
	return js
}

type PlayerObject struct {
	UserObject
	Heros []string `json:"heros" bson:"heros"`
}
