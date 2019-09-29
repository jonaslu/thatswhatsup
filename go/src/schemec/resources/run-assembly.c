#include <stdio.h>
#include <wchar.h>
#include <locale.h>
#include <stdint.h>

#include <inttypes.h>

extern uint64_t scheme_entry();

#define int_tag 0b00
#define int_bitmask 0b00000011

#define boolean_tag 0b0011111
#define boolean_mask 0b1111111

#define char_tag 0b00001111
#define char_mask 0b11111111

#define empty_list 0b00101111

#define pair_tag 0b00000001
#define pair_mask 0b00000111

int printPrimitive(uint64_t result)
{
  if ((result & int_bitmask) == int_tag)
  {
    result = result >> 2;
    printf("%d", result);

    return 0;
  }
  else if ((result & boolean_mask) == boolean_tag)
  {
    result = result >> 7;
    if (result == 1)
    {
      printf("true");
      return 0;
    }

    printf("false");
    return 1;
  }
  else if (result == empty_list)
  {
    printf("()");
    return 0;
  }
  else if ((result & char_mask) == char_tag)
  {
    result = result >> 8;

    setlocale(LC_CTYPE, "");
    wchar_t test = result;

    wprintf(L"#\\%lc", test);

    return 0;
  }

  return 1;
}

int printHeapValue(uint64_t result)
{
  if ((result & pair_mask) == pair_tag)
  {
    int64_t *car_addr = (uint64_t *)(result - 1);
    int64_t *cdr_addr = (uint64_t *)(result + 7);

    printf("(");
    if (printHeapValue(*car_addr)) {
      printf("Unknown value type %d", result);
      return 1;
    }
    printf(" ");

    printHeapValue(*cdr_addr); // Need to recurse here too

    printf(")");
  }
  else
  {
    if (printPrimitive(result) != 0)
    {
      printf("Unknown value type %d", result);
      return 1;
    }
  }

  return 0;
}

int main(void)
{
  uint64_t result = scheme_entry();
  return printHeapValue(result);
}
