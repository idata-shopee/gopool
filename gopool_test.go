package gopool

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, expect interface{}, actual interface{}, message string) {
	if expect == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expect %v !=  actual %v", expect, actual)
	}
	t.Fatal(message)
}

func TestBase(t *testing.T) {
	type Res = func(int, int) int
	res := func(a int, b int) int {
		return a + b
	}

	getNewItem := func(onItemBoken OnItemBorken) (Item, error) {
		return Item{res, func() {}}, nil
	}

	pool := GetPool(getNewItem, 8)

	for i := 1; i < 1000; i++ {
		api, _ := pool.Get()
		assertEqual(t, api.(Res)(i, 1), i+1, "")
		assertEqual(t, api.(Res)(i, 2), i+2, "")
		assertEqual(t, api.(Res)(i, 3), i+3, "")
	}
}
