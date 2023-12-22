package hy

import (
	"sync"
)

var defaultPool = Pool{pool: sync.Pool{New: func() any {
	b := make([]byte, 2048)
	return b
},
},
	maxCount:  10000,
	currCount: 0,
}

type Pool struct {
	pool      sync.Pool
	maxCount  int64
	currCount int64
}

func (p *Pool) Get() []byte {
	//for {
	//	if atomic.LoadInt64(&p.currCount) <= p.maxCount {
	//		break
	//	} else {
	//		time.Sleep(time.Millisecond * 20)
	//	}
	//}
	//atomic.AddInt64(&p.currCount, 1)
	b := p.pool.Get()
	return b.([]byte)
}

func (p *Pool) Put(b []byte) {
	p.pool.Put(b)
	//atomic.AddInt64(&p.currCount, -1)
}
