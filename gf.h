#ifndef GF_H
#define GF_H

#include <stdint.h>

uint8_t gf_add(uint8_t lhs, uint8_t rhs);
uint8_t gf_sub(uint8_t lhs, uint8_t rhs);
uint8_t gf_mul(uint8_t lhs, uint8_t rhs);
uint8_t gf_div(uint8_t dividend, uint8_t divisor);

unsigned int gf_log(uint8_t v);
uint8_t gf_alog(unsigned int v);
uint8_t gf_log_mul(unsigned int lhs, uint8_t rhs);

#endif
