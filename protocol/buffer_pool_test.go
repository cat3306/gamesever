package protocol

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

func TestInitBufferPool(t *testing.T) {
	//InitBufferPool()
	list := make([]int, 0)
	for key := range BUFFERPOOL.pool {
		list = append(list, int(key))
	}
	sort.Ints(list)
	if len(list) != len(BUFFERPOOL.capSlice) {
		t.Fatal("err 1")
	}
	for i := 0; i < len(BUFFERPOOL.capSlice); i++ {
		if BUFFERPOOL.capSlice[i] != uint32(list[i]) {
			t.Fatal("err 2")
		}
	}

}
func TestBufferPool_GetBuffer(t *testing.T) {
	//l := len(BUFFERPOOL.capSlice)
	var i uint32
	for i = 1; i < 100000; i++ {
		buff := BUFFERPOOL.Get(i)
		if math.Log2(float64(len(buff))) != math.Log2(float64(i))+7 {
			fmt.Println(math.Log2(float64(len(buff))), math.Ceil(math.Log2(float64(i))))
			t.Fatalf("11 need:%d,get:%d", i, len(buff))
		}
		t.Logf("need:%d,get:%d", i, len(buff))
		BUFFERPOOL.Put(buff)
	}

}
func TestBufferPool(t *testing.T) {
	//l := len(BUFFERPOOL.capSlice)
	var i uint32
	for i = 1; i < 100000; i++ {
		buff := make([]int, 100)
		t.Logf("need:%d,get:%d", i, len(buff))
	}
}
