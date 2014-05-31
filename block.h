#ifndef BLOCK_H
#define BLOCK_H

#include <stdint.h>

#include "encode.h"

void block_encode(codec_t *codec, uint8_t *block, int block_len);
void block_decode(codec_t *codec, uint8_t *block, int block_len, uint8_t *chunks);

#endif
