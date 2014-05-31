package greeso

/*
 #include "gf.h"
*/
import "C"

func add(lhs, rhs byte) byte {
	return byte(C.gf_add(C.uint8_t(lhs), C.uint8_t(rhs)))
}

func sub(lhs, rhs byte) byte {
	return byte(C.gf_sub(C.uint8_t(lhs), C.uint8_t(rhs)))
}

func mul(lhs, rhs byte) byte {
	return byte(C.gf_mul(C.uint8_t(lhs), C.uint8_t(rhs)))
}

func div(dividend, divisor byte) byte {
	return byte(C.gf_div(C.uint8_t(dividend), C.uint8_t(divisor)))
}
