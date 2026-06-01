// Package pool provides a type-safe generic wrapper around sync.Pool.
package pool

import (
	"sync"
)

// Pool is a type-safe wrapper around sync.Pool with custom generation and resetting.
type Pool[T any] struct {
	pool *sync.Pool

	reset func(*T) *T
}

// Option configures a Pool.
type Option[T any] func(*Pool[T])

// New returns an Option that sets the generator function for creating new objects.
func New[T any](new func() *T) Option[T] {
	return func(p *Pool[T]) {
		p.pool = &sync.Pool{
			New: func() any { return new() },
		}
	}
}

// Reset returns an Option that sets the reset function run on every object retrieved via Get.
func Reset[T any](reset func(*T) *T) Option[T] {
	return func(p *Pool[T]) {
		p.reset = reset
	}
}

// NewPool creates a Pool with the given options.
// If no options are specified, it defaults to using new(T) and performs no resetting.
func NewPool[T any](opts ...Option[T]) Pool[T] {
	p := Pool[T]{
		reset: func(t *T) *T { return t },
	}

	for _, o := range opts {
		o(&p)
	}

	if p.pool == nil {
		p.pool = &sync.Pool{
			New: func() any { return new(T) },
		}
	}

	return p
}

// Get retrieves an object from the Pool, resets it, and returns it.
// If the pool is empty, it uses the generator function or new(T) to create a new object.
func (p Pool[T]) Get() *T {
	return p.reset(p.pool.Get().(*T))
}

// Put returns an object to the Pool for future reuse.
func (p Pool[T]) Put(t *T) {
	p.pool.Put(t)
}
