package test

import (
	//"encoding/json"
	"fmt"
	//"github.com/bitly/go-simplejson"
	"mzlib"
	"testing"
)

func TestUserModel(t *testing.T) {
	fmt.Println("UserModel test begin.......")

	um := mzlib.UserModelGet("15618973660")
	fmt.Printf("UserModel um ToJson=======%s\n", um.ToJson())
	um.Put()

	um2 := mzlib.UserModelGet("15618973660")
	fmt.Printf("UserModel um2 ToJson=======%s\n", um2.ToJson())

	// um0 := new(mzlib.UserModel)
	// um0.Pk = "Uid"
	// um0.Uid = "15618973660"
	// um0.Info = map[string]interface{}{
	// 	"ggg": "uuuu",
	// }
	// um0.Name = "xushicai"
	// um0.Put()

	// um := new(mzlib.UserModel)
	// um.Pk = "Uid"

	// um2 := um.Get("15618973660").(*mzlib.UserModel)
	// fmt.Printf("UserModel um2 ToJson=======%s\n", um2.ToJson())
	// fmt.Printf("UserModel um2.Uid=======%s\n", um2.Uid)
	// fmt.Printf("UserModel um2 ToBson=======%s\n", um2.ToBson())
	// um2.Name = "ppp"
	// um2.DoPut()
	// um2.Delete()

	// type UserCards struct {
	// 	mzlib.UserModel `bson:",inline"`
	// 	Cards           []string `json:"cards" bson:"cards"`
	// }

	// uc := new(UserCards)
	// uc.Pk = "Uid"
	// uc.Uid = "166666666"
	// uc.Name = "xushicai"
	// uc.Info = map[string]interface{}{
	// 	"yyyy": "xxxx",
	// }
	// uc.Cards = []string{"card_1", "card_2"}
	// uc.Put()
	// uc.DoPut()

	fmt.Println("UserModel test end.........")
}
