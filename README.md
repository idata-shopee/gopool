# gopool

A pool library for golang

## Quick Example

```go
import (
  "github.com/lock-free/gopool"
  "time"
)

// define resource
type Res = func(int, int) int
res := func(a int, b int) int {
	return a + b
}

// create pool
getNewItem := func(onItemBoken gopool.OnItemBorken) (*gopool.Item, error) {
	return &gopool.Item{res, func() {}}, nil
}

pool := GetPool(getNewItem, 8, 3000*time.Millisecond)

// Get resource

res, err := pool.Get()

if err != nil {
  // handle
} else {
  res.(Res)(2, 3) // 5
}
```
