#include <inttypes.h>
#include <stdio.h>

#include "headers.h"
#include "utils.h"

#define ARP_HEADER_SIZE 28

void handle_arp_header(unsigned char *payload) {
  struct arp_hdr *hdr = (struct arp_hdr *)payload;

  dump_as_hex(payload, ARP_HEADER_SIZE);

  printf("Arp hw type: %X\n", hdr->hwtype);
}
