package greeso

/*
 #include "matrix.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"strconv"
)

type (
	vector []byte
	row    []uint
	matrix []row
)

var (
	ErrDimensionMismatch = errors.New("dimension mismatch")
	ErrNonInvert         = errors.New("matrix not invertible")
)

// NewMatrix initializes a zero matrix of m rows and n columns.
func NewMatrix(m, n int) matrix {
	A := make([]row, m)
	for i := range A {
		A[i] = make([]uint, n)
	}
	return A
}

// Copy the matrix.
func (A matrix) Copy() matrix {
	m := NewMatrix(len(A), len(A[0]))
	for i, row := range A {
		copy(m[i], row)
	}
	return m
}

// Identity generates a (possibly truncated) identity matrix.
func Identity(m, n int) matrix {
	I := NewMatrix(m, n)
	for diag := range I {
		if diag < len(I[diag]) {
			I[diag][diag] = uint(1)
		}
	}
	return I
}

func (A matrix) Mul(V vector, C vector) error {
	if len(A[0]) != len(V) {
		return ErrDimensionMismatch
	}
	for i := range C {
		C[i] = 0
	}
	for i := range A {
		for j, a := range A[i] {
			C[i] ^= mul(byte(a), V[j])
		}
	}
	return nil
}

func (A matrix) Inverse() matrix {
	w := A.Copy()
	result := Identity(len(A), len(A[0]))
	w.LowerGaussianElim(result)
	w.UpperInverse(result)
	return result
}

func (A matrix) PartialLowerGaussElim(row, col int, inverse matrix) (int, int) {
	lastRow := len(A) - 1
	for row < lastRow {
		if col >= len(A[row]) {
			return row, col
		}
		if 0 == A[row][col] {
			return row, col
		}
		divisor := div(byte(1), byte(A[row][col]))
		for k := row + 1; k < len(A); k++ {
			nextTerm := A[k][col]
			if nextTerm == 0 {
				continue
			}
			multiple := mul(divisor, sub(byte(0), byte(nextTerm)))
			A.rowMulAdd(multiple, row, k)
			if inverse != nil {
				inverse.rowMulAdd(multiple, row, k)
			}
		}
		row = row + 1
		col = col + 1
	}
	return row, col
}

func (A matrix) LowerGaussianElim(inverse matrix) matrix {
	row, col := 0, 0
	lastRow, lastCol := len(A)-1, len(A[0])-1
	if lastRow > lastCol+1 {
		lastRow = lastCol + 1
	}
	for row < lastRow && col < lastCol {
		leader := A.findRowLeader(row, col)
		if leader < 0 {
			col = col + 1
			continue
		}
		if leader != row {
			A.rowAdd(leader, row)
			if inverse != nil {
				inverse.rowAdd(leader, row)
			}
		}
		row, col = A.PartialLowerGaussElim(row, col, inverse)
	}
	return A
}

func (A matrix) UpperInverse(inverse matrix) (matrix, error) {
	lastCol := len(A)
	if lastCol > len(A[0]) {
		lastCol = len(A[0])
	}
	for col := 0; col < lastCol; col++ {
		if byte(0) == byte(A[col][col]) {
			return nil, ErrNonInvert
		}
		divisor := div(byte(1), byte(A[col][col]))
		if divisor != byte(1) {
			A.rowMul(col, divisor, 0)
			if inverse != nil {
				inverse.rowMul(col, divisor, 0)
			}
		}
		for elim := 0; elim < col; elim++ {
			multiple := sub(byte(0), byte(A[elim][col]))
			A.rowMulAdd(multiple, col, elim)
			if inverse != nil {
				inverse.rowMulAdd(multiple, col, elim)
			}
		}
	}
	return A, nil
}

func (A matrix) Transpose() matrix {
	old := A
	A = make([]row, len(old[0]))
	for row := range A {
		A[row] = make([]uint, len(old))
	}
	for r := range A {
		for c := range old {
			A[r][c] = old[c][r]
		}
	}
	return A
}

func (A matrix) Logify() matrix {
	old := A
	A = make([]row, len(old))
	for row := range A {
		A[row] = make([]uint, len(old[row]))
	}
	for i, row := range old {
		for j, v := range row {
			A[i][j] = log(byte(v))
		}
	}
	return A
}

func (A matrix) AntiLogify() matrix {
	old := A
	A = make([]row, len(old))
	for row := range A {
		A[row] = make([]uint, len(old[row]))
	}
	for i, row := range old {
		for j, v := range row {
			A[i][j] = uint(aLog(v))
		}
	}
	return A
}

func (A matrix) LogMul(V vector, C vector) error {
	if len(A[0]) != len(V) {
		return ErrDimensionMismatch
	}
	for i := range C {
		C[i] = 0
	}
	for i := range A {
		for j, a := range A[i] {
			C[i] ^= logMul(a, V[j])
		}
	}
	return nil
}

func (A matrix) String() string {
	return A.GoString()
}

func (A matrix) GoString() string {
	m := 0
	for _, row := range A {
		for _, c := range row {
			l := len(strconv.Itoa(int(c)))
			if l > m {
				m = l
			}
		}
	}
	s := ""
	f := "%" + strconv.Itoa(m+1) + "s"
	for _, r := range A {
		s = s + "\n"
		for _, c := range r {
			s = s + fmt.Sprintf(f, strconv.Itoa(int(c)))
		}
	}
	return s
}

func (A matrix) rowMul(row int, multiplier byte, start int) {
	for i := range A[row] {
		A[row][i] = uint(mul(byte(A[row][i]), multiplier))
	}
}

func (A matrix) rowAdd(i, j int) {
	for k := range A[j] {
		A[j][k] = uint(add(byte(A[j][k]), byte(A[i][k])))
	}
}

func (A matrix) rowMulAdd(multiplier byte, i, j int) {
	for k := range A[j] {
		A[j][k] = uint(add(byte(A[j][k]), mul(multiplier, byte(A[i][k]))))
	}
}

func (A matrix) findRowLeader(row, col int) int {
	for r := row; r < len(A); r++ {
		if uint(0) != A[r][col] {
			return r
		}
	}
	return -1
}
