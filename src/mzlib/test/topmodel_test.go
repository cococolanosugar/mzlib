package test

import (
	"fmt"
	"mzlib"
	"testing"
)

func TestTopModel(t *testing.T) {
	fmt.Println("TopModel test begin.......")
	tm := mzlib.NewTopModel("bocai_gold_rank")

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("name%d", i)
		score := i
		tm.Set(name, score)
	}

	list := tm.Get(10, false)
	fmt.Println("list=======", list)

	desc_list := tm.Get(10, true)
	fmt.Println("desc_list=======", desc_list)
	for i := 0; i < len(list); i++ {
		fmt.Println("list[i]====", list[i], "======desc_list[len-i-1]======", desc_list[len(list)-i-1])
		if list[i] != desc_list[len(list)-i-1] {
			panic("assert fail!!!")
		}
	}

	tm.Remove("name9")
	cnt := tm.Count()
	if cnt != 9 {
		panic("assert fail!")
	}

	tm.RemoveByScore(7, 8)
	cnt = tm.Count()
	if cnt != 7 {
		panic("assert fail!")
	}

	tm.RemoveLast(2, false)
	cnt = tm.Count()
	if cnt != 5 {
		panic("assert fail!")
	}

	tm.RemoveLast(2, true)
	cnt = tm.Count()
	if cnt != 3 {
		panic("assert fail!")
	}

	rank := tm.Rank("name2", false)
	if rank != 0 {
		panic("assert fail!")
	}
	rank = tm.Rank("name2", true)
	if rank != 2 {
		panic("assert fail!")
	}

	rank = tm.Rank("meiyou", true)
	if rank != -1 {
		panic("assert fail!")
	}

	score := tm.Score("name3")
	if score != float64(3) {
		panic("assert fail!")
	}

	score = tm.Score("meiyou")
	if score != float64(-1) {
		panic("assert fail!")
	}

	fmt.Println("TopModel test end.........")
}
