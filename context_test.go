package mygin

import (
	"encoding/json"
	"sync"
	"testing"
)

var contextPool = sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c := &Context{}

		json.Marshal(c)
	}
}

func BenchmarkUnmarshalWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c := contextPool.Get().(*Context)
		json.Marshal(c)
		contextPool.Put(c)
	}
}
