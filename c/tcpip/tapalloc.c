#include <fcntl.h>
#include <net/if.h>
#include <linux/if.h>
#include <linux/if_tun.h>
#include <sys/ioctl.h>

#include <stdio.h> // printf
#include <stdlib.h> // exits
#include <string.h> // memset
#include <unistd.h> // read


int tun_alloc()
{
  struct ifreq ifr;
  int fd, err;

  memset(&(ifr), 0, sizeof(ifr));

  int fd_open_code = fd = open("/dev/net/tun", O_RDWR);
  if (fd_open_code < 0) {
    printf("Error opening device, code %d", fd_open_code);
    exit(1);
  }

  ifr.ifr_flags = IFF_TAP | IFF_NO_PI;

  int ioctl_code = ioctl(fd, TUNSETIFF, (void *) &ifr);
  if (ioctl_code < 0) {
    printf("Error setting ioctl flags, code %d", ioctl_code);
    exit(1);
  }

  return fd;
}

void write_ether_frame(int fd) {

}


void read_ether_frame(int fd) {
  char data[32];
  int bytes_read = read(fd, &data, 32);
  printf("%d bytes read, msg %s\n", bytes_read, data);
}

int main(void) {
  int fd = tun_alloc();

  while (1) {
    read_ether_frame(fd);
  }

  return 0;
}

// Read bytes from the network and dump out what happens
