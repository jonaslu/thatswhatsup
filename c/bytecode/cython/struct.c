#include <stdio.h>

typedef struct Object {
  int i;
  // char* str;
} Object;

void print_test(Object* test) {
  printf("Getting in here\n");
  printf("Value of test %p", test);
  printf("I val: %d\n", test->i);
  // printf("Char val: %s\n", test->test);
}