#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

#include "encode.h"
#include "matrix.h"

extern codec_t* codec_new(int n, int k) {
  codec_t *new;
  new = malloc(sizeof *new);
  if (new == NULL) {
	 return NULL;
  }

  new->encode = matrix_new(n, k);
  new->decode = matrix_new(k, k);

  // Make a vandermonde matrix.
  for (int i=0; i < n; ++i) {
	 uint8_t term = 1;
	 for (int j=0; j < k; ++j) {
		matrix_set(new->encode, i, j, term);
		term = gf_mul(term, i);
	 }
  }

  // Make it a systemic encoder by attempting to solve the transpose
  // (the upper n x n will be the identity matrix).
  matrix_transpose(new->encode);
  matrix_lower_gauss(new->encode, NULL);
  matrix_upper_inverse(new->encode, NULL);
  matrix_transpose(new->encode);

  return new;
}

extern void codec_free(codec_t *c) {
  if (c == NULL) return;
  matrix_free(c->encode);
  matrix_free(c->decode);
  free(c);
}

extern void codec_encode(codec_t *c, uint8_t *message, uint8_t *code) {
  vector_t m_v = {c->encode->n};
  vector_t c_v = {c->encode->m};
  m_v.d = message;
  c_v.d = code;
  matrix_mul(c->encode, &m_v, &c_v);
}

extern void codec_prepare_decoder(codec_t *c, uint8_t* chunks) {
  for (int i=0; i < c->decode->n; ++i) {
	 matrix_copy_row(c->decode, i, c->encode, chunks[i]);
  }
  matrix_inverse(c->decode);
}

extern void codec_decode(codec_t *c, uint8_t *message, uint8_t *code) {
  vector_t m_v = {c->decode->n};
  vector_t c_v = {c->decode->n};
  m_v.d = message;
  c_v.d = code;
  matrix_mul(c->decode, &c_v, &m_v);
}
