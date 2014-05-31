#ifndef ENCODE_H
#define ENCODE_H

#include <stdint.h>

#include "matrix.h"

typedef struct {
  matrix_t *encode;
  matrix_t *decode;
} codec_t;

codec_t* codec_new(int n, int k);
void codec_free(codec_t *c);
void codec_encode(codec_t *c, uint8_t *message, uint8_t *code);
void codec_prepare_decoder(codec_t *c, uint8_t* chunks);
void codec_decode(codec_t *c, uint8_t *message, uint8_t *code);

#endif
