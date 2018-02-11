#include <stdio.h>
#include <wchar.h>
#include <locale.h>

extern int scheme_entry();

#define int_tag 0b00
#define int_bitmask 0b00000011

#define boolean_tag 0b0011111
#define boolean_mask 0b1111111

#define char_tag 0b00001111
#define char_mask 0b11111111

#define empty_list 0b00101111

int main(void) {
  int result = scheme_entry();


  if ((result & int_bitmask) == int_tag) {
    result = result >> 2;
    printf("%d", result);
  } else if ((result & boolean_mask) == boolean_tag) {
    result = result >> 8;
    if (result == 1) {
      printf("true");
    } else {
      printf("false");
    }
  } else if (result == empty_list) {
    printf("()");
  } else if ((result & char_mask) == char_tag) {
    result = result >> 8;

    setlocale(LC_CTYPE, "");
    wchar_t test = result;

    wprintf(L"#\\%lc", test);
  } else {
    printf("Unknown value type %d", result);
  }

  return 0;
}
