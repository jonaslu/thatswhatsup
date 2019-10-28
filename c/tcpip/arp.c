#include <inttypes.h>
#include <stdio.h>
#include <net/if_arp.h>

#include "headers.h"
#include "utils.h"

void handle_arp_header(unsigned char *payload) {
  struct arp_hdr *hdr = (struct arp_hdr *)payload;

  if (hdr->hwtype == ARPHRD_ETHER && hdr->prosize&& hdr->opcode == ARPOP_REQUEST) {

  }

  dump_arp_hdr(payload);
}
