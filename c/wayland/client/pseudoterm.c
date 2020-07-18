#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <pty.h>
#include <sys/select.h>
#include <errno.h>

#include "client.h"

#define BUF_SIZE 10000

char output[BUF_SIZE];
int masterfd;

int init_pty()
{
  int pid = forkpty(&masterfd, NULL, NULL, NULL);

  if (pid < 0)
  {
    printf("Could not fork pty\n");
    exit(1);
  }

  if (pid == 0)
  {
    char *argv[1];
    if (execvp("/bin/bash", argv) == -1)
    {
      printf("Failed to start bash");
      exit(1);
    };

    printf("Returning after execvp");
  }

  return masterfd;
}

void read_from_pty(struct wl_client *wl_client)
{
  char buffer[BUF_SIZE];

  if (read(masterfd, buffer, BUF_SIZE) == -1)
  {
    printf("Read error from master\n");
    printf("Error: %s\n", strerror(errno));
    exit(1);
  }

  strcat(output, buffer);

  render_chars(wl_client, output);
}

void write_to_pty(const char *text, int size)
{
  if (write(masterfd, text, size) != size)
  {
    printf("Error writing to master\n");
    exit(1);
  };
}
