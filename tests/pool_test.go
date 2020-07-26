package tests

import (
	"fmt"
	"jim/common/tool"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	wait := sync.WaitGroup{}
	defer tool.ReleaseGoPool()
	list := []object{
		{index: 2, name: "abc"},
		{index: 3, name: "erw"},
		{index: 1, name: "vds"},
		{index: 4, name: "hku"},
		{index: 5, name: "pta"},
	}
	for _, obj := range list {
		o := obj
		wait.Add(1)
		executor := func() {
			time.Sleep(time.Duration(o.index) * time.Second)
			run(o)
			wait.Done()
		}
		tool.AsyncRun(executor)
	}
	wait.Wait()
	fmt.Println("complected")
}

func run(obj object) {
	fmt.Println(obj.index, obj.name)
}

type object struct {
	index int
	name  string
}
