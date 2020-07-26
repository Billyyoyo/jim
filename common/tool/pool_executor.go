package tool

import (
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	// go程池大小
	DefaultAntsPoolSize = 1 << 8

	// go程有效期
	ExpiryDuration = 60 * time.Second

	// 是否非阻塞
	Nonblocking = true

	// 是否预分配内存
	PreAlloc = true
)

var (
	pool *ants.Pool
)

func init() {
	options := ants.Options{
		ExpiryDuration: ExpiryDuration,
		PreAlloc:       PreAlloc,
		Nonblocking:    Nonblocking}
	var err error
	pool, err = ants.NewPool(DefaultAntsPoolSize, ants.WithOptions(options))
	if err != nil {
		log.Error("create goroutine pool fail:", err)
		panic("create goroutine pool fail: " + err.Error())
	}
}

func AsyncRun(f func()){
	pool.Submit(f)
}

func ReleaseGoPool(){
	pool.Release()
}