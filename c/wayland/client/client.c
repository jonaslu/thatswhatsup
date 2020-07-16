#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <limits.h>

#include <unistd.h>
#include <pty.h>
#include <poll.h>

#include "wayland-client.h"
#include "client.h"

static void init_display(struct wl_client *client)
{
  client->display = wl_display_connect(NULL);
  if (!client->display)
  {
    fprintf(stderr, "Failed to connect to Wayland display.\n");
    exit(1);
  }

  fprintf(stderr, "Connection established!\n");
}

int main()
{
  struct wl_client wl_client;

  init_display(&wl_client);
  init_registry_listener(&wl_client);

  wl_display_dispatch(wl_client.display);
  wl_display_roundtrip(wl_client.display);

  assert(wl_client.display && wl_client.compositor && wl_client.shm && wl_client.wm_base && wl_client.seat);

  init_surface(&wl_client);

  wl_display_flush(wl_client.display);
  wl_display_dispatch(wl_client.display);

  int waylandFd = wl_display_get_fd(wl_client.display);
  int ptyFd = init_pty();

  struct pollfd pollfd[] = {{
                                .fd = waylandFd,
                                .events = POLLIN,
                            },
                            {
                                .fd = ptyFd,
                                .events = POLLIN,
                            }};

  // I want the fd for the pty & wayland. So I can poll the pty fd here?

  for (;;)
  {
    wl_display_flush(wl_client.display);
    if (poll(pollfd, 1, INT_MAX) == -1)
    {
      printf("Poll error\n");
      exit(1);
    }

    if (pollfd[0].revents |= POLLIN)
    {
      wl_display_dispatch(wl_client.display);
    }

    if (pollfd[1].revents |= POLLIN)
    {
      read_from_pty(&wl_client);
    }
  }

  wl_display_roundtrip(wl_client.display);
  wl_display_disconnect(wl_client.display);

  return 0;
}
