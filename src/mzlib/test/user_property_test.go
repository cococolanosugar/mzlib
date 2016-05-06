package test

import (
	//"encoding/json"
	"fmt"
	//"github.com/bitly/go-simplejson"
	"mzlib"
	"testing"
)

func TestUserProperty(t *testing.T) {
	fmt.Println("UserProperty test begin.......")

	up := mzlib.NewUserProperty()
	up.Uid = "123456"
	up.Name = "PaPaPa"
	up.Title = "CTO"
	up.Gold = 100
	up.Coin = 1000
	up.Put()

	fmt.Println("up ToJson===========", up.ToJson())
	fmt.Println("up ToBson===========", up.ToBson())

	up2 := mzlib.UserPropertyGet("123456")
	fmt.Println("up2 ToJson==========", up2.ToJson())

	fmt.Println("UserProperty test end.........")
}
