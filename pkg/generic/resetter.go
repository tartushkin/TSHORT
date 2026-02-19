package generic

import (
	"sync"
)

// Resetter — интерфейс, требующий наличия метода Reset.
type Resetter interface {
	Reset()
}

type Pool[T any] struct {
	pool sync.Pool
}

func New[T Resetter]() *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				var zero T
				return zero
			},
		},
	}
}

func (p *Pool[T]) Get() T {
	v := p.pool.Get()
	if v == nil {
		var zero T
		return zero
	}
	return v.(T)
}

func (p *Pool[T]) Put(obj T) {
	// obj должен быть таким, чтобы можно было вызвать Reset()
	resetter, ok := interface{}(obj).(Resetter)
	if ok {
		resetter.Reset()
	}
	p.pool.Put(obj)
}
