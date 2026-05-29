// Package helper provides pre-configured pools and a global registry for sharing them.
package helper

import (
	"bytes"
	"reflect"
	"sync"

	"github.com/shu-go/pool"
)

// NewBufferPool creates a Pool of bytes.Buffer objects.
// Created buffers are initialized with the given capacity and are reset when retrieved.
func NewBufferPool(capacity int) pool.Pool[bytes.Buffer] {
	return pool.NewPool(
		pool.New(func() *bytes.Buffer {
			obj := &bytes.Buffer{}
			obj.Grow(capacity)
			return obj
		}),
		pool.Reset(func(obj *bytes.Buffer) *bytes.Buffer {
			obj.Reset()
			return obj
		}))
}

// NewBytesPool creates a Pool of []byte slice pointers.
// Created slices are initialized with the given capacity and are reset to length zero when retrieved.
func NewBytesPool(capacity int) pool.Pool[[]byte] {
	return pool.NewPool(
		pool.New(func() *[]byte {
			obj := make([]byte, 0, capacity)
			return &obj
		}),
		pool.Reset(func(obj *[]byte) *[]byte {
			*obj = (*obj)[:0]
			return obj
		}))
}

var (
	registry = sync.Map{}
)

func init() {
	RegisterPool(NewBufferPool(8))
	RegisterPool(NewBytesPool(8))
}

// RegisterPool registers a Pool of the type T in the global registry.
func RegisterPool[T any](pool pool.Pool[T]) {
	registry.Store(reflect.TypeFor[T](), pool)
}

// NewPoolOf retrieves a registered Pool of the type T from the global registry.
// If not registered, it returns a default Pool using new(T) and a no-op reset function.
func NewPoolOf[T any]() pool.Pool[T] {
	if p, found := registry.Load(reflect.TypeFor[T]()); found {
		return p.(pool.Pool[T])
	}

	return pool.NewPool(
		pool.New(func() *T { return new(T) }),
		pool.Reset(func(t *T) *T { return t }),
	)
}
