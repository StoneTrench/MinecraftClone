package async

import (
	"fmt"
	"sync"
)

type Future[T any] struct {
	done chan struct{}
	once sync.Once
	val  T
	err  error
}

func newFuture[T any]() *Future[T] {
	return &Future[T]{done: make(chan struct{})}
}

func (f *Future[T]) set(v T, err error) {
	f.once.Do(func() {
		f.val, f.err = v, err
		close(f.done)
	})
}

func (f *Future[T]) Await() (T, error) {
	<-f.done
	return f.val, f.err
}

func (f *Future[T]) AwaitPanic() T {
	<-f.done
	if f.err != nil {
		panic(fmt.Errorf("error occured in async.Future: %w", f.err))
	}
	return f.val
}

func (f *Future[T]) Done() <-chan struct{} { return f.done }

func NewFn[T any](fn func() (T, error)) *Future[T] {
	f := newFuture[T]()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				var zero T
				f.set(zero, fmt.Errorf("future panicked: %v", r))
			}
		}()
		v, err := fn()
		f.set(v, err)
	}()
	return f
}
func NewSet[T any]() (*Future[T], func(T, error)) {
	f := newFuture[T]()
	return f, f.set
}

func Must[T any](f *Future[T]) T {
	v, err := f.Await()
	if err != nil {
		panic(err)
	}
	return v
}
