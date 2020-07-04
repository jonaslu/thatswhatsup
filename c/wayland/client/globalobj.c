#include <stdio.h>
#include <string.h>

#include "wayland-client.h"
#include "xdg-shell-client-protocol.h"

#include "client.h"

static void global_registry_handler(void *data, struct wl_registry *wl_registry, uint32_t id, const char *interface, uint32_t version)
{
  struct wl_client *wl_client = data;

  printf("Got an event for %s id %d\n", interface, id);

  if (strcmp(interface, wl_compositor_interface.name) == 0)
  {
    wl_client->compositor = wl_registry_bind(wl_registry, id, &wl_compositor_interface, version);
  }

  if (strcmp(interface, xdg_wm_base_interface.name) == 0)
  {
    wl_client->wm_base = wl_registry_bind(wl_registry, id, &xdg_wm_base_interface, version);
  }

  if (strcmp(interface, wl_shm_interface.name) == 0)
  {
    wl_client->shm = wl_registry_bind(wl_registry, id, &wl_shm_interface, version);
  }
}

static void global_remove(void *our_data,
                          struct wl_registry *registry,
                          uint32_t name)
{
  printf("Got remove: %d", name);
}

static struct wl_registry_listener registry_listener = {
    .global = global_registry_handler,
    .global_remove = global_remove,
};

void init_registry_listener(struct wl_client *client)
{
  struct wl_registry *registry = wl_display_get_registry(client->display);
  wl_registry_add_listener(registry, &registry_listener, client);
}
