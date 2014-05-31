#include <stdint.h>
#include <stdlib.h>

#include "block.h"
#include "matrix.h"

void block_encode(codec_t *codec, uint8_t *block, int block_len) {
  vector_t m = {codec->encode->n};
  vector_t c = {codec->encode->m};
  uint8_t m_buf[256];
  uint8_t c_buf[256];
  m.d = m_buf;
  c.d = c_buf;

  int stripe_len = block_len / m.n;
  for (int i=0; i < stripe_len; ++i) {
	 for (int j=0; j < m.n; ++j) {
		m.d[j] = block[i + j * stripe_len];
	 }
	 codec_encode(codec, m.d, c.d);
	 for (int j = m.n; j < c.n; ++j) {
		block[i + j * stripe_len] = c.d[j];
	 }
  }
}

void block_decode(codec_t *codec, uint8_t *block, int block_len, uint8_t *chunks) {
  codec_prepare_decoder(codec, chunks);
  vector_t m = {codec->decode->n};
  vector_t c = {codec->decode->m};
  uint8_t m_buf[256];
  uint8_t c_buf[256];
  m.d = m_buf;
  c.d = c_buf;

  int stripe_len = block_len / m.n;
  for (int i=0; i < stripe_len; ++i) {
	 for (int j=0; j < c.n; ++j) {
		c.d[j] = block[i + j * stripe_len];
	 }
	 codec_decode(codec, m.d, c.d);
	 for (int j = 0; j < m.n; ++j) {
		block[i + j * stripe_len] = m.d[j];
	 }
  }
}
