#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include "wayland-client.h"

#include "client.h"

static void init_display(struct wl_client *client) {
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

  assert(wl_client.display && wl_client.compositor && wl_client.shm && wl_client.wm_base);

  add_and_render_surface(&wl_client);

  while (wl_display_dispatch(wl_client.display))
  {
  };

  wl_display_roundtrip(wl_client.display);
  wl_display_disconnect(wl_client.display);

  return 0;
}
