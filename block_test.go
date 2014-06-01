package greeso

import (
	"os"
	"runtime/pprof"
	"testing"
)

func TestEncodeBlock(t *testing.T) {
	m := []byte("Hello World!")
	l := len(m)
	chunker := NewCBlockCodec(8, 5)
	b := chunker.Encode(m)
	t.Log(string(b))
	b = b[3:18]
	t.Log(string(b))
	chunks := []byte{1, 2, 3, 4, 5}
	d := chunker.Decode(b, chunks)[0:l]
	t.Log(string(d))
	if string(m) != string(d) {
		t.Error(string(d))
	}
}

func BenchmarkCEncodeBlock(b *testing.B) {
	m := make([]byte, 16000, 32000)
	chunker := NewCBlockCodec(8, 5)
	for i := 0; i < b.N; i++ {
		chunker.Encode(m)
	}
}

func BenchmarkEncodeBlock(b *testing.B) {
	m := make([]byte, 16000, 32000)
	chunker := NewBlockCodec(8, 5)
	f, _ := os.Create("encodeblock.prof")
	b.ResetTimer()
	pprof.StartCPUProfile(f)
	for i := 0; i < b.N; i++ {
		chunker.Encode(m)
	}
	pprof.StopCPUProfile()
}
