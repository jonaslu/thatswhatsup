#include <stdio.h>

extern int scheme_entry();

#define int_bitmask 0b00000011
#define int_tag 0b00

#define rest_bitmask 0b11111111

#define char_tag 0b00001111
#define boolean_tag 0b0011111
#define empty_list 0b00101111

int main(void) {
  int result = scheme_entry();

  if ((result & int_bitmask) == int_tag) {
    result = result >> 2;
    printf("%d", result);
  } else {
    printf("Unknown value type %d", result);
  }

  return 0;
}
