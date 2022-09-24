package example

import (
	"fmt"
	"github.com/valyala/bytebufferpool"
	"testing"
)

func TestBufferPool(t *testing.T) {
	bytebufferpool.Put(&bytebufferpool.ByteBuffer{
		B: make([]byte, 1024),
	})
	bytebufferpool.Put(&bytebufferpool.ByteBuffer{
		B: make([]byte, 2048),
	})
	p := bytebufferpool.Pool{}
	p.Get()
	for i := 0; i < 100000; i++ {

		fmt.Println(cap(bytebufferpool.Get().B))
	}
}

