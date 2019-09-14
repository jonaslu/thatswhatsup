#include <stdio.h>

#include "headers.h"

void dump_as_hex(unsigned char *buf, int bytes_read)
{
  for (int i = 0; i < bytes_read; i++)
  {
    printf("%02X ", buf[i]);
  }
}

void dump_eth_hdr(unsigned char *hdr, int bytes_read) {
  int position = 0;
  printf("dmac: ");

  for (int i = 1; i<6; i++) {
    printf("%02X:", hdr[position++]);
  }
  printf("%02X\n",hdr[position++]);

  printf("smac: ");
  for (int i=1; i<6; i++) {
    printf("%02X:", hdr[position++]);
  }
  printf("%02X\n",hdr[position++]);

  printf("ethertype: 0x%02X%02X\n", hdr[position++], hdr[position++]);
}
