package pool_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/shu-go/gotwant"

	"github.com/shu-go/pool"
	"github.com/shu-go/pool/helper"
)

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

func TestBytesPool(t *testing.T) {
	p := helper.NewBytesPool(10)
	b := p.Get()

	gotwant.Test(t, len(*b), 0)

	*b = append(*b, "hoge"...)
	gotwant.Test(t, string(*b), "hoge")
	gotwant.TestExpr(t, cap(*b), cap(*b) >= 4)

	bb := b
	p.Put(b)
	b = p.Get()
	gotwant.Test(t, len(*b), 0)
	gotwant.TestExpr(t, cap(*b), cap(*b) >= 4)
	gotwant.Test(t, b, bb)
	p.Put(b)
}

func TestBuffer(t *testing.T) {
	p := helper.NewBufferPool(8)
	b := p.Get()

	gotwant.Test(t, b.Len(), 0)

	b.WriteString("hoge")
	gotwant.Test(t, b.String(), "hoge")
	gotwant.TestExpr(t, b.Cap(), b.Cap() >= 4)

	bb := b
	p.Put(b)
	b = p.Get()
	gotwant.Test(t, b.Len(), 0)
	gotwant.TestExpr(t, b.Cap(), b.Cap() >= 4)
	gotwant.Test(t, b, bb)
	p.Put(b)
}

func TestRegistered(t *testing.T) {
	p1 := helper.NewPoolOf[bytes.Buffer]()
	buf := p1.Get()
	gotwant.Test(t, buf.Len(), 0)
	gotwant.TestExpr(t, buf.Cap(), buf.Cap() >= 8)

	p2 := helper.NewPoolOf[[]byte]()
	s := p2.Get()
	gotwant.Test(t, len(*s), 0)
	gotwant.TestExpr(t, cap(*s), cap(*s) >= 8)

}

func BenchmarkBuffer(b *testing.B) {
	b.Run("Manual", func(b *testing.B) {
		p := sync.Pool{
			New: func() any {
				b := &bytes.Buffer{}
				b.Grow(10)
				return b
			},
		}
		b.ResetTimer()
		for b.Loop() {
			b := p.Get().(*bytes.Buffer)
			b.Reset()
			b.WriteString("hello")
			p.Put(b)
		}
	})

	b.Run("Pool", func(b *testing.B) {
		p := helper.NewBufferPool(8)
		b.ResetTimer()
		for b.Loop() {
			b := p.Get()
			b.WriteString("hello")
			p.Put(b)
		}
	})

	b.Run("Pool(reg)", func(b *testing.B) {
		p := helper.NewPoolOf[bytes.Buffer]()
		b.ResetTimer()
		for b.Loop() {
			b := p.Get()
			b.WriteString("hello")
			p.Put(b)
		}
	})
}

func BenchmarkBytes(b *testing.B) {
	b.Run("Manual", func(b *testing.B) {
		p := sync.Pool{
			New: func() any {
				s := make([]byte, 0, 10)
				return &s
			},
		}
		b.ResetTimer()
		for b.Loop() {
			s := p.Get().(*[]byte)
			*s = (*s)[:0]
			*s = append(*s, "hello"...)
			p.Put(s)
		}
	})

	b.Run("Pool", func(b *testing.B) {
		p := helper.NewBytesPool(10)
		b.ResetTimer()
		for b.Loop() {
			s := p.Get()
			*s = append(*s, "hello"...)
			p.Put(s)
		}
	})

	b.Run("Pool(reg)", func(b *testing.B) {
		p := helper.NewPoolOf[[]byte]()
		b.ResetTimer()
		for b.Loop() {
			s := p.Get()
			*s = append(*s, "hello"...)
			p.Put(s)
		}
	})
}

func BenchmarkIntReg(b *testing.B) {
	p := helper.NewPoolOf[int]()
	for b.Loop() {
		i := p.Get()
		*i = 123
		p.Put(i)
	}
}
