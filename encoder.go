package greeso

/*
 #cgo CFLAGS: -O3
 #include "encode.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type matrixCodec struct {
	encode matrix
	decode matrix
}

type cCodec struct {
	codec *C.codec_t
}

func NewCCodec(n, k int) *cCodec {
	c := new(cCodec)
	c.codec = C.codec_new(C.int(n), C.int(k))
	runtime.SetFinalizer(c, func(c *cCodec) {
		C.codec_free(c.codec)
	})
	return c
}

func (c *cCodec) Encode(message, code []byte) []byte {
	C.codec_encode(c.codec,
		(*C.uint8_t)(unsafe.Pointer(&message[0])),
		(*C.uint8_t)(unsafe.Pointer(&code[0])))
	return code
}

func (c *cCodec) PrepareDecoder(chunks []byte) {
	C.codec_prepare_decoder(c.codec, (*C.uint8_t)(unsafe.Pointer(&chunks[0])))
}

func (c *cCodec) Decode(code []byte, message []byte) []byte {
	C.codec_decode(c.codec,
		(*C.uint8_t)(unsafe.Pointer(&message[0])),
		(*C.uint8_t)(unsafe.Pointer(&code[0])))
	return message
}

func NewRSCodec(n, k int) matrixCodec {
	m := matrixCodec{}
	encode := NewMatrix(n, k)
	for i := 0; i < n; i++ {
		term := byte(1)
		for j := 0; j < k; j++ {
			encode[i][j] = term
			term = mul(term, byte(i))
		}
	}
	encode = encode.Transpose()
	encode = encode.LowerGaussianElim(nil)
	encode, _ = encode.UpperInverse(nil)
	encode = encode.Transpose()
	m.encode = encode
	return m
}

func (m *matrixCodec) Encode(message []byte, code []byte) []byte {
	m.encode.Mul(message, code)
	return code
}

func (m *matrixCodec) PrepareDecoder(chunks []byte) {
	decode := NewMatrix(len(chunks), len(chunks))
	for i, r := range chunks {
		decode[i] = m.encode[r]
	}
	m.decode = decode.Inverse()
}

func (m *matrixCodec) Decode(code []byte, message []byte) []byte {
	m.decode.Mul(code, message)
	return message
}
