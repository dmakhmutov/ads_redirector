package round_robin

import (
	"sync"
)

type RoundRobin struct {
	sync.Mutex

	current int
	pool    []string
}

func Stub() *RoundRobin {
	return &RoundRobin{
		current: 0,
		pool:    []string{"http://ya.ru"},
	}
}

func (r *RoundRobin) New(pool []string) {
	r.Lock()

	defer r.Unlock()
	r.pool = pool
}

func (r *RoundRobin) Next() string {
	r.Lock()
	defer r.Unlock()

	if r.current >= len(r.pool) {
		r.current = r.current % len(r.pool)
	}

	result := r.pool[r.current]
	r.current++
	return result
}
