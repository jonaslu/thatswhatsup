#include <inttypes.h>
#include <stdio.h>
#include <net/if_arp.h>
#include <linux/if_ether.h>
#include <arpa/inet.h>

#include "headers.h"
#include "utils.h"

void handle_arp_header(unsigned char *payload)
{
  struct arp_hdr *hdr = (struct arp_hdr *)payload;

  dump_arp_hdr(payload);

  if (hdr->hwtype == htons(ARPHRD_ETHER) && hdr->protype == htons(ETH_P_IP))
  {
    dump_arp_ipv4(hdr->data, hdr->hwsize, hdr->prosize);

    struct arp_ipv4 *payload = (struct arp_ipv4 *)hdr->data;
  }
}
