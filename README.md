Package pool provides a type-safe generic wrapper around sync.Pool.

[![Go Report Card](https://goreportcard.com/badge/github.com/shu-go/pool)](https://goreportcard.com/report/github.com/shu-go/pool)
[![Go Reference](https://pkg.go.dev/badge/github.com/shu-go/pool.svg)](https://pkg.go.dev/github.com/shu-go/pool)
![MIT License](https://img.shields.io/badge/License-MIT-blue)

# go get

```
go get -u github.com/shu-go/pool
```

# Example

```go
func Example_pool() {
	var p pool.Pool[bytes.Buffer]
	p = pool.NewPool(
		pool.New(func() *bytes.Buffer {
			return &bytes.Buffer{}
		}),
		pool.Reset(func(t *bytes.Buffer) *bytes.Buffer {
			t.Reset()
			return t
		}))

	var buf *bytes.Buffer
	buf = p.Get()
	p.Put(buf)
}

func Example_helper() {
	var p pool.Pool[bytes.Buffer]
	p = helper.NewBufferPool(16)

	var buf *bytes.Buffer
	buf = p.Get()
	p.Put(buf)
}

func Example_helperFor() {
	var p pool.Pool[bytes.Buffer]
	p = helper.NewPoolOf[bytes.Buffer]()

	var buf *bytes.Buffer
	buf = p.Get()
	p.Put(buf)
}
```
