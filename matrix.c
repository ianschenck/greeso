#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "gf.h"
#include "matrix.h"

// Internal forward declarations.
uint8_t* at(matrix_t *m, int row, int column);
void matrix_row_mul(matrix_t *m, int row, uint8_t multiple);
void matrix_row_add(matrix_t *m, int i, int j);
void matrix_row_muladd(matrix_t *m, uint8_t multiple, int i, int j);
int matrix_row_leader(matrix_t *m, int row, int col);
point_t matrix_part_lower_gauss(matrix_t *m, point_t p, matrix_t *inverse);

extern uint8_t* at(matrix_t *m, int i, int j) {
  return &m->d[m->n * i + j];
}

extern matrix_t* matrix_new(int m, int n) {
  matrix_t *matrix;
  matrix = calloc(sizeof *matrix + m * n, 1);
  matrix->m = m;
  matrix->n = n;
  return matrix;
}

extern void matrix_free(matrix_t *matrix) {
  if (matrix == NULL) return;
  free(matrix);
}

extern void matrix_copy(matrix_t *dst, matrix_t *src) {
  size_t m_size = sizeof *src + src->m*src->n;
  memcpy(dst, src, m_size);
}

extern matrix_t* matrix_duplicate(matrix_t *matrix) {
  size_t m_size = sizeof *matrix + matrix->m*matrix->n;
  matrix_t *new = malloc(m_size);
  memcpy(new, matrix, m_size);
  return new;
}

extern void matrix_copy_row(matrix_t *dst, int d_row, matrix_t *src, int s_row) {
  size_t row_size = sizeof(uint8_t) * src->n;
  memcpy(dst->d + d_row * row_size, src->d + s_row * row_size, row_size);
}

extern void matrix_zero(matrix_t *matrix) {
  memset(matrix->d, 0, matrix->m * matrix->n * sizeof(uint8_t));
}

extern void matrix_identity(matrix_t *m) {
  for (int i=0; i < m->m && i < m->n; ++i) {
	 *at(m, i, i) = 1;
  }
}

extern int matrix_mul(matrix_t *m, vector_t *v, vector_t *c) {
  if (m->n != v->n) {
	 return 0;
  }

  for (int i=0; i < m->m; ++i) {
		c->d[i] = 0;
	 for (int j=0; j < m->n; ++j) {
		c->d[i] = c->d[i] ^ gf_mul(*at(m, i, j), v->d[j]);
	 }
  }
  return -1;
}

extern void matrix_inverse(matrix_t *m) {
  matrix_t *result = matrix_new(m->m, m->n);

  matrix_identity(result);
  matrix_lower_gauss(m, result);
  matrix_upper_inverse(m, result);
  matrix_copy(m, result);

  matrix_free(result);
}

extern void matrix_lower_gauss(matrix_t *m, matrix_t *inverse) {
  point_t p = {0,0};
  int last_row = m->m - 1;
  int last_col = m->n - 1;
  if (last_row > m->n) {
	 last_row = m->n;
  }
  while (p.i < last_row && p.j < last_col) {
	 int leader = matrix_row_leader(m, p.i, p.j);
	 if (leader < 0) {
		p.j += 1;
		continue;
	 }
	 if (leader != p.i) {
		matrix_row_add(m, leader, p.i);
		if (inverse != NULL) {
		  matrix_row_add(inverse, leader, p.i);
		}
	 }
	 p = matrix_part_lower_gauss(m, p, inverse);
  }
}

extern int matrix_upper_inverse(matrix_t *m, matrix_t *inverse) {
  int last_col = m->m < m->n ? m->m : m->n;
  for (int j=0; j < last_col; ++j) {
	 if (*at(m, j, j) == 0) {
		return 0;
	 }
	 uint8_t divisor = gf_div(1, *at(m, j, j));
	 if (divisor != 1) {
		matrix_row_mul(m, j, divisor);
		if (inverse != NULL) {
		  matrix_row_mul(inverse, j, divisor);
		}
	 }
	 for (int elim=0; elim < j; ++elim) {
		uint8_t multiple = gf_sub(0, *at(m, elim, j));
		matrix_row_muladd(m, multiple, j, elim);
		if (inverse != NULL) {
		  matrix_row_muladd(inverse, multiple, j, elim);
		}
	 }
  }
  return -1;
}

// Transposition in-place is less clear, so do it the inefficient way.
extern void matrix_transpose(matrix_t *m) {
  matrix_t *old = matrix_duplicate(m);

  m->m = old->n;
  m->n = old->m;

  for (int i=0; i < m->m; ++i) {
	 for (int j=0; j < m->n; ++j) {
		uint8_t v = matrix_get(old, j, i);
		matrix_set(m, i, j, v);
	 }
  }

  matrix_free(old);
}

point_t matrix_part_lower_gauss(matrix_t *m, point_t p, matrix_t *inverse) {
  int last_row = m->m - 1;
  for (int i=0; i < last_row; ++i) {
	 if (p.j >= m->n) {
		return p;
	 }
	 if (*at(m, p.i, p.j) == 0) {
		return p;
	 }
	 uint8_t divisor = gf_div(1, *at(m, p.i, p.j));
	 for (int k=p.i+1; k < m->m; ++k) {
		uint8_t next_term = *at(m, k, p.j);
		if (next_term == 0) {
		  continue;
		}
		uint8_t multiple = gf_mul(divisor, gf_sub(0, next_term));
		matrix_row_muladd(m, multiple, p.i, k);
		if (inverse != NULL) {
		  matrix_row_muladd(inverse, multiple, p.i, k);
		}
	 }
	 p.i += 1;
	 p.j += 1;
  }
  return p;
}

void matrix_row_mul(matrix_t *m, int row, uint8_t multiple) {
  for (int j=0; j < m->n; ++j) {
	 *at(m, row, j) = gf_mul(*at(m, row, j), multiple);
  }
}

void matrix_row_add(matrix_t *m, int i, int j) {
  for (int k=0; k < m->n; ++k) {
	 *at(m, j, k) = gf_add(*at(m, j, k), *at(m, i, k));
  }
}

void matrix_row_muladd(matrix_t *m, uint8_t multiple, int i, int j) {
  for (int k=0; k < m->n; ++k) {
	 *at(m, j, k) = gf_add(*at(m, j, k), gf_mul(multiple, *at(m, i, k)));
  }
}

int matrix_row_leader(matrix_t *m, int row, int col) {
  for (int r=row; r < m->m; ++r) {
	 if (*at(m, r, col) != 0) {
		return r;
	 }
  }
  return -1;
}
