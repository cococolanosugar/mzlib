package test

import (
	"fmt"
	"mzlib"
	"testing"
)

func TestSequence(t *testing.T) {
	fmt.Println("Sequence test begin.......")
	s := mzlib.NextUid()
	fmt.Println("s===============", s)
	fmt.Println("Sequence test end.........")
}
