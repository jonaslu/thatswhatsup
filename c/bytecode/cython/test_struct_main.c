#include "struct.h"

int main(void) {
  Object test;
  test.i = 2;
  test.test = "yeye";

  print_test(&test);

  return 0;
}