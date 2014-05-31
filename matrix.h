#ifndef MATRIX_H
#define MATRIX_H

#include <stdint.h>

#include "gf.h"

typedef struct {
  int n;
  uint8_t *d;
} vector_t;

typedef struct {
  int m;
  int n;
  uint8_t d[];
} matrix_t;

typedef struct {
  int i;
  int j;
} point_t;

matrix_t* matrix_new(int m, int n);
void matrix_free(matrix_t *matrix);
matrix_t* matrix_duplicate(matrix_t *matrix);
void matrix_zero(matrix_t *matrix);
void matrix_identity(matrix_t *matrix);
void matrix_copy_row(matrix_t *dst, int d_row, matrix_t *src, int s_row);
int matrix_mul(matrix_t *m, vector_t *v, vector_t *c);
void matrix_log(matrix_t *m);
int matrix_log_mul(matrix_t *m, vector_t *v, vector_t *c);
void matrix_inverse(matrix_t *m);
void matrix_transpose(matrix_t *m);
void matrix_lower_gauss(matrix_t *m, matrix_t *inverse);
int matrix_upper_inverse(matrix_t *m, matrix_t *inverse);

inline uint8_t matrix_get(matrix_t *m, int i, int j) {
  return m->d[m->n * i + j];
}

inline void matrix_set(matrix_t *m, int i, int j, uint8_t x) {
  m->d[m->n * i + j] = x;
}

#endif
