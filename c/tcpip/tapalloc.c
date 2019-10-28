#include <fcntl.h>
#include <net/if.h>
#include <linux/if.h>
#include <linux/if_tun.h>
#include <linux/if_ether.h>
#include <sys/ioctl.h>
#include <arpa/inet.h>

#include <stdio.h>    // printf
#include <stdlib.h>   // exits
#include <string.h>   // memset
#include <unistd.h>   // read

#include "headers.h"
#include "arp.h"
#include "utils.h"

struct eth_hdr *get_ether_header(char *buf)
{
  struct eth_hdr *hdr = (struct eth_hdr *)buf;

  // It's big endian in the linux/if_ether.h
  // This turns the host (le) to network (be)
  hdr->ethertype = htons(hdr->ethertype);

  return hdr;
}

int tun_alloc()
{
  struct ifreq ifr;
  int fd, err;

  memset(&(ifr), 0, sizeof(ifr));

  int fd_open_code = fd = open("/dev/net/tun", O_RDWR);
  if (fd_open_code < 0)
  {
    printf("Error opening device, code %d", fd_open_code);
    exit(1);
  }

  ifr.ifr_flags = IFF_TAP | IFF_NO_PI;

  int ioctl_code = ioctl(fd, TUNSETIFF, (void *)&ifr);
  if (ioctl_code < 0)
  {
    printf("Error setting ioctl flags, code %d", ioctl_code);
    exit(1);
  }

  return fd;
}

int main(void)
{
  int buf_len = 32;
  char buf[32];
  int fd = tun_alloc();

  while (1)
  {
    int bytes_read = read(fd, &buf, buf_len);
    dump_eth_hdr(buf);
    struct eth_hdr *hdr = get_ether_header(buf);

    switch (hdr->ethertype)
    {
    case ETH_P_ARP:
    {
      handle_arp_header(hdr->payload);
      break;
    }
    default:
    {
      break;
    }
    }
  }

  return 0;
}

// Read bytes from the network and dump out what happens
