package greeso

import (
	"testing"
)

func TestEncoder(t *testing.T) {
	message := "Hello"
	code := make([]byte, 8)
	m := NewRSCodec(8, 5)
	m.Encode([]byte(message), code)
	t.Logf("encoded '%s'", string(code))
	partialCode := append(code[0:1], code[3], code[5], code[6], code[7])
	t.Logf("decoding '%s'", string(partialCode))
	chunks := []byte{0, 3, 5, 6, 7}
	m.PrepareDecoder(chunks)
	recovered := make([]byte, 5)
	m.Decode(partialCode, recovered)
	t.Logf("decoded '%s'", string(recovered))
	if message != string(recovered) {
		t.Errorf("'%s' != '%s'", recovered, message)
	}
}

func BenchmarkEncoder(b *testing.B) {
	message := "Hello"
	code := make([]byte, 8)
	m := NewRSCodec(8, 5)
	for i := 0; i < b.N; i++ {
		m.Encode([]byte(message), code)
	}
}

func TestCCodec(t *testing.T) {
	message := "Hello"
	code := make([]byte, 8)
	m := NewCCodec(8, 5)
	m.Encode([]byte(message), code)
	t.Logf("encoded '%s'", string(code))
	partialCode := append(code[0:1], code[3], code[5], code[6], code[7])
	t.Logf("decoding '%s'", string(partialCode))
	chunks := []byte{0, 3, 5, 6, 7}
	m.PrepareDecoder(chunks)
	recovered := make([]byte, 5)
	m.Decode(partialCode, recovered)
	t.Logf("decoded '%s'", string(recovered))
	if message != string(recovered) {
		t.Errorf("'%s' != '%s'", recovered, message)
	}
}
