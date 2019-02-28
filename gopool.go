package gopool

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
	"sync"
	"time"
)

// usage of pool:
//   1. backup
//   2. dynamic LB (TODO)

type CleanFunction = func()

type Item struct {
	Resouce interface{}   // keep Resouce
	Clean   CleanFunction // Clean Resouce TODO with DLB
}

type OnItemBorken = func()

// get new item and can know the moment when it brokes
type GetNewItem = func(OnItemBorken) (*Item, error)

// define pool data structure
type Pool struct {
	items      map[string]*Item
	getNewItem GetNewItem
	size       int
	mutex      *sync.Mutex
	duration   time.Duration
	status     int // pool status, 0: stoped, 1: opened
}

func (pool *Pool) GetItemNum() int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	return len(pool.items)
}

func (pool *Pool) addNewItem() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.items) < pool.size {
		uid, uerr := uuid.NewV4()
		id := uid.String()
		if uerr != nil {
			// TODO
			fmt.Println(uerr)
		} else {
			item, err := pool.getNewItem(func() {
				// when item broken, remove from pool
				go pool.removeItem(id)
			})
			if err != nil {
				// TODO
				fmt.Println(err)
			} else {
				pool.items[id] = item
			}
		}
	}
}

func (pool *Pool) removeItem(id string) {
	pool.mutex.Lock()
	delete(pool.items, id)
	pool.mutex.Unlock()

	pool.addNewItem()
}

func (pool *Pool) Get() (interface{}, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	poolLen := len(pool.items)
	if poolLen <= 0 {
		return 0, errors.New("pool is empty")
	}

	// random
	source := rand.NewSource(time.Now().UnixNano())
	ran := rand.New(source)
	n := ran.Intn(poolLen)

	count := 0
	for key := range pool.items {
		if count == n {
			return pool.items[key].Resouce, nil
		}
		count += 1
	}
	return 0, errors.New("pool is empty!!!")
}

// shut down a pool
// normally this is used for testing
func (pool *Pool) Shutdown() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	pool.status = 0 // change status

	// clean all resources
	for id, item := range pool.items {
		item.Clean()
		delete(pool.items, id)
	}
}

func (pool *Pool) maintain() {
	if pool.status == 1 { // only maintain alive pool
		pool.addNewItem()
		go (func() {
			time.Sleep(pool.duration)
			// keep maintain
			pool.maintain()
		})()
	}
}

// get a pool, need to provide:
//    (1) getNewItem: how to get a new item
//    (2) size
//    (3) duration to get a new item
func GetPool(getNewItem GetNewItem, size int, duration time.Duration) Pool {
	items := map[string]*Item{}
	pool := Pool{items, getNewItem, size, &sync.Mutex{}, duration, 1}
	pool.maintain()

	return pool
}
