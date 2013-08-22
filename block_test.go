package greeso

import (
	"testing"
)

func TestEncodeBlock(t *testing.T) {
	m := []byte("Hello World!")
	l := len(m)
	chunker := NewBlockCodec(8, 5)
	b := chunker.Encode(m)
	b = b[3:18]
	chunks := []int{1, 2, 3, 4, 5}
	if string(m) != string(chunker.Decode(b, chunks)[0:l]) {
		t.Error(string(chunker.Decode(b, chunks)))
	}
}

func BenchmarkEncodeBlock(b *testing.B) {
	m := make([]byte, 5000, 8000)
	chunker := NewBlockCodec(8, 5)
	for i := 0; i < b.N; i++ {
		chunker.Encode(m)
	}
}
