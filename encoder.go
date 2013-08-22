package greeso

import ()

type matrixCodec struct {
	encode matrix
	decode matrix
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
	encode = encode.Transpose().LowerGaussianElim(nil)
	encode, _ = encode.UpperInverse(nil)
	encode = encode.Transpose()
	m.encode = encode
	return m
}

func (m *matrixCodec) Encode(message []byte, code []byte) []byte {
	m.encode.Mul(message, code)
	return code
}

func (m *matrixCodec) PrepareDecoder(chunks []int) {
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
