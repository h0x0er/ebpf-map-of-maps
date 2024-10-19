#include "common.h"

struct inner_map {
  __uint(type, BPF_MAP_TYPE_LRU_HASH);
  __type(key, u32);
  __type(value, u8);
  __uint(max_entries, 1);
} inner_map SEC(".maps");

struct {
  __uint(type, BPF_MAP_TYPE_HASH_OF_MAPS);
  __uint(max_entries, 1024);
  __type(key, u32);
  __array(values, struct inner_map);
} outer_map SEC(".maps");