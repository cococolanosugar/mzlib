package test

import (
	"fmt"
	"mzlib"
	"testing"
	"time"
)

func TestTmpModel(t *testing.T) {
	fmt.Println("TmpModel test begin.......")
	key := "tmpmodelusecase"
	value := "gogogogo"
	obj := mzlib.TmpModelSet(key, value, time.Hour)

	s := mzlib.TmpModelGet(key)
	if obj.Value != s {
		panic("assert fail!!!")
	}

	fmt.Println("TmpModel test end.........")
}
