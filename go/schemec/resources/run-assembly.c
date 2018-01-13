#include <stdio.h>

extern int scheme_entry();

#define int_bitmask 3
#define int_tag 0

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
