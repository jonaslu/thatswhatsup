#include <stdio.h>
#include <string.h>

#include "xdg-shell-client-protocol.h"
#include "wayland-client.h"

static struct wl_client
{
  struct wl_display *display;
  struct wl_compositor *compositor;
  struct wl_registry *registry;
  struct xdg_wm_base *wm_base;
} wl_client;

static void global_registry_handler(void *data, struct wl_registry *wl_registry, uint32_t id, const char *interface, uint32_t version)
{
  struct wl_client *wl_client = data;

  printf("Got an event for %s id %d\n", interface, id);

  if (strcmp(interface, "wl_compositor") == 0)
  {
    wl_client->compositor = wl_registry_bind(wl_registry, id, &wl_compositor_interface, version);
  }

  if (strcmp(interface, "xdg_wm_base") == 0)
  {
    wl_client->wm_base = wl_registry_bind(wl_registry, id, &xdg_wm_base_interface, version);
  }
}

static void global_remove(void *our_data,
                          struct wl_registry *registry,
                          uint32_t name)
{
  printf("Got remove: %d", name);
}

struct wl_registry_listener registry_listener = {
    .global = global_registry_handler,
    .global_remove = global_remove};

int main()
{
  struct wl_client wl_client;

  wl_client.display = wl_display_connect(NULL);
  if (!wl_client.display)
  {
    fprintf(stderr, "Failed to connect to Wayland display.\n");
    return 1;
  }

  fprintf(stderr, "Connection established!\n");

  wl_client.registry = wl_display_get_registry(wl_client.display);
  wl_registry_add_listener(wl_client.registry, &registry_listener, &wl_client);

  wl_display_roundtrip(wl_client.display);

  wl_display_disconnect(wl_client.display);
  return 0;
}
