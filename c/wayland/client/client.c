#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/syscall.h>
#include <sys/mman.h>

#include "xdg-shell-client-protocol.h"
#include "wayland-client.h"

static struct wl_client
{
  struct wl_display *display;
  struct wl_compositor *compositor;
  struct wl_registry *registry;
  struct xdg_wm_base *wm_base;
  struct wl_shm *shm;
  struct wl_surface *surface;
} wl_client;

const int width = 200, height = 200;
int stride = 4;
int size = 200 * 200 * 4;

static int allocate_buffer(unsigned int **data)
{
  int fd = syscall(SYS_memfd_create, "buffer", 0);
  ftruncate(fd, size);

  *data = mmap(NULL, size, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
  if (data == MAP_FAILED)
  {
    printf("Failed to allocate shared buffer\n");
    assert(NULL);
  }

  return fd;
}

static struct wl_buffer *render_buffer(struct wl_shm *shm)
{
  unsigned int *data;
  int fd = allocate_buffer(&data);

  struct wl_shm_pool *pool = wl_shm_create_pool(shm, fd, size);
  struct wl_buffer *buffer = wl_shm_pool_create_buffer(pool, 0, width, height, width * stride, WL_SHM_FORMAT_XRGB8888);
  wl_shm_pool_destroy(pool);
  close(fd);

  for (int y = 0; y < height; ++y)
  {
    for (int x = 0; x < width; ++x)
    {
      if ((x + y / 8 * 8) % 16 < 8)
        data[y * width + x] = 0xFF666666;
      else
        data[y * width + x] = 0xFFEEEEEE;
    }
  }

  munmap(data, size);

  return buffer;
}

static void xdg_surface_configure(void *data, struct xdg_surface *surface, uint32_t serial)
{
  struct wl_client *client = data;

  xdg_surface_ack_configure(surface, serial);
  struct wl_buffer *buffer = render_buffer(client->shm);
  wl_surface_attach(client->surface, buffer, 0, 0);
  wl_surface_commit(client->surface);
}

const static struct xdg_surface_listener surface_listener = {
    .configure = xdg_surface_configure,
};

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

struct wl_registry_listener registry_listener = {
    .global = global_registry_handler,
    .global_remove = global_remove,
};

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

  wl_display_dispatch(wl_client.display);
  wl_display_roundtrip(wl_client.display);

  assert(wl_client.display && wl_client.compositor && wl_client.shm && wl_client.wm_base);

  wl_client.surface = wl_compositor_create_surface(wl_client.compositor);
  struct xdg_surface *xdg_surface = xdg_wm_base_get_xdg_surface(wl_client.wm_base, wl_client.surface);

  struct xdg_toplevel *xdg_toplevel = xdg_surface_get_toplevel(xdg_surface);
  xdg_toplevel_set_title(xdg_toplevel, "My f example");

  xdg_surface_add_listener(xdg_surface, &surface_listener, &wl_client);
  wl_surface_commit(wl_client.surface);

  while (wl_display_dispatch(wl_client.display))
  {
  };

  wl_display_roundtrip(wl_client.display);
  wl_display_disconnect(wl_client.display);

  return 0;
}
