package async

import (
	"fmt"
	"sync"
)

type Future[T any] struct {
	done chan struct{}
	val  T
	err  error
	once sync.Once
}

func (f *Future[T]) TryAwait() (T, error) {
	<-f.done
	return f.val, f.err
}

func (f *Future[T]) AwaitPanic() T {
	v, err := f.TryAwait()
	if err != nil {
		panic(err)
	}
	return v
}

func NewFn[T any](fn func() (T, error)) *Future[T] {
	f := &Future[T]{
		done: make(chan struct{}),
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.err = fmt.Errorf("future panicked: %v", r)
			}
			close(f.done)
		}()
		f.val, f.err = fn()
	}()

	return f
}

func NewSet[T any]() (*Future[T], func(T, error)) {
	f := &Future[T]{
		done: make(chan struct{}),
	}

	return f, func(t T, err error) {
		f.once.Do(func() {
			f.val = t
			f.err = err
			close(f.done)
		})
	}
}
