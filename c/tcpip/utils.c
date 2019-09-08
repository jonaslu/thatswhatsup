#include <stdio.h>

void dump_as_hex(char *buf, int bytes_read)
{
  for (int i = 0; i < bytes_read; i++)
  {
    printf("%02X ", buf[i]);
  }
}
