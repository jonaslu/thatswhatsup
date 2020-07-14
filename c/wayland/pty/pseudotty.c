#define _POSIX_C_SOURCE 200112L
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <pty.h>
#include <sys/select.h>

#define BUF_SIZE 256

int masterfd;
pid_t pid;

void main()
{
  pid = forkpty(&masterfd, NULL, NULL, NULL);

  if (pid < 0)
  {
    printf("Could not fork pty\n");
    exit(1);
  }

  if (pid == 0)
  {
    // Child
    execvp("/bin/bash", (char **)NULL);
  }
  else
  {
    for (;;)
    {
      char buf[BUF_SIZE];
      fd_set read_fd, write_fd, error_fd;

      FD_ZERO(&read_fd);
      FD_ZERO(&write_fd);
      FD_ZERO(&error_fd);

      FD_SET(STDIN_FILENO, &read_fd);
      FD_SET(masterfd, &read_fd);

      if (select(masterfd + 1, &read_fd, NULL, NULL, NULL) == -1)
      {
        printf("Error doing select\n");
        exit(1);
      }

      if (FD_ISSET(STDIN_FILENO, &read_fd))
      {
        printf("Reading from stdin\n");
        int numread = read(STDIN_FILENO, buf, BUF_SIZE);
        if (numread < 0)
        {
          printf("Could not read from STDIN\n");
          // exit(1);
        }

        if (write(masterfd, buf, numread) != numread)
        {
          printf("Write to masterfd failed\n");
          // exit(1);
        }
      }

      if (FD_ISSET(masterfd, &read_fd))
      {
        printf("Reading from masterfd\n");
        int numread = read(masterfd, buf, BUF_SIZE);
        if (numread < 0)
        {
          printf("Could not read from masterfd\n");
          // exit(1);
        }

        if (write(STDOUT_FILENO, buf, numread) != numread)
        {
          printf("Write to STDOUT failed\n");
          // exit(1);
        }
      }
    }
  }
}
