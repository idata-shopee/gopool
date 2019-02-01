package gopool

import (
	"errors"
	"fmt"
	"testing"
	"time"
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

// simple case
func TestBase(t *testing.T) {
	type Res = func(int, int) int
	res := func(a int, b int) int {
		return a + b
	}

	getNewItem := func(onItemBoken OnItemBorken) (*Item, error) {
		return &Item{res, func() {}}, nil
	}

	pool := GetPool(getNewItem, 8, 3000*time.Millisecond)

	for i := 1; i < 1000; i++ {
		api, _ := pool.Get()
		assertEqual(t, api.(Res)(i, 1), i+1, "")
		assertEqual(t, api.(Res)(i, 2), i+2, "")
		assertEqual(t, api.(Res)(i, 3), i+3, "")
	}
}

// on item broken case
func TestOnBroken(t *testing.T) {
	type Res = func(int, int) int
	res := func(a int, b int) int {
		return a + b
	}

	count := 0
	getNewItem := func(onItemBoken OnItemBorken) (*Item, error) {
		count += 1
		if count == 3 {
			go (func() {
				time.Sleep(50 * time.Millisecond)
				onItemBoken()
			})()
		}
		return &Item{res, func() {}}, nil
	}

	pool := GetPool(getNewItem, 3, 50*time.Millisecond)

	time.Sleep(250 * time.Millisecond)
	assertEqual(t, 3, pool.GetItemNum(), "")
}

func TestOnBroken2(t *testing.T) {
	type Res = func(int, int) int
	res := func(a int, b int) int {
		return a + b
	}

	count := 0
	getNewItem := func(onItemBoken OnItemBorken) (*Item, error) {
		count += 1
		if count == 3 {
			go (func() {
				time.Sleep(50 * time.Millisecond)
				onItemBoken()
			})()
		}

		if count > 3 {
			return nil, errors.New("no items")
		}
		return &Item{res, func() {}}, nil
	}

	pool := GetPool(getNewItem, 3, 1*time.Millisecond)

	time.Sleep(40 * time.Millisecond)
	assertEqual(t, 3, pool.GetItemNum(), "")
	time.Sleep(200 * time.Millisecond)
	assertEqual(t, 2, pool.GetItemNum(), "")
}
