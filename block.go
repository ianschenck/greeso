package greeso

import (
	"math"
)

type block []byte

type blockCodec struct {
	n int
	k int
	matrixCodec
}

func NewBlockCodec(n, k int) *blockCodec {
	return &blockCodec{n, k, NewRSCodec(n, k)}
}

func (c *blockCodec) Encode(block []byte) []byte {
	message := make([]byte, c.k, c.n)
	code := make([]byte, c.n)
	stripeLen := int(math.Ceil(float64(len(block)) / float64(c.k)))
	encodedLen := stripeLen * c.n
	if cap(block) < encodedLen {
		tempBlock := make([]byte, len(block), encodedLen)
		copy(tempBlock, block)
		block = tempBlock
	}
	block = block[0:encodedLen]
	for i := 0; i < stripeLen; i++ {
		for j := 0; j < c.k; j++ {
			offset := i + j*stripeLen
			message[j] = block[offset]
		}
		c.matrixCodec.Encode(message, code)
		for j := c.k; j < c.n; j++ {
			offset := i + j*stripeLen
			block[offset] = code[j]
		}
	}
	return block
}

func (c *blockCodec) Decode(block []byte, chunks []int) []byte {
	code := make([]byte, c.k)
	message := make([]byte, c.k)
	stripeLen := int(math.Ceil(float64(len(block)) / float64(c.k)))
	decodedLen := stripeLen * c.k
	if cap(block) < decodedLen {
		tempBlock := make([]byte, len(block), decodedLen)
		copy(tempBlock, block)
		block = tempBlock
	}
	block = block[0:decodedLen]
	c.PrepareDecoder(chunks)
	for i := 0; i < stripeLen; i++ {
		for j := 0; j < c.k; j++ {
			offset := i + j*stripeLen
			code[j] = block[offset]
		}
		c.matrixCodec.Decode(code, message)
		for j := 0; j < c.k; j++ {
			offset := i + j*stripeLen
			block[offset] = message[j]
		}
	}
	return block
}
