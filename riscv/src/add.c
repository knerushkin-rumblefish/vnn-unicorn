#include <stdint.h>
#define INPUT_LEN_ADDRESS 0x40100000
#define INPUT_ADDRESS 0x40100008
#define OUTPUT_ADDRESS 0x40000000

int main() {
  int64_t *input_len_ptr = (int64_t *)INPUT_LEN_ADDRESS;
  int32_t *input = (int32_t *)INPUT_ADDRESS;

  int64_t blob_size = *input_len_ptr;

  int64_t blob_acc = 0;

  for (int64_t i = 0; i < blob_size; i += 4) {
    blob_acc = blob_acc + *input;
    input++;
  };

  int64_t *output = (int64_t *)OUTPUT_ADDRESS;
  *output = blob_acc;

  return 0;
}
