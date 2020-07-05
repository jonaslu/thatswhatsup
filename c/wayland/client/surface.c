#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/syscall.h>
#include <sys/mman.h>

#include "xdg-shell-client-protocol.h"

#include "client.h"

int add = 0;

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
    exit(1);
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
        data[y * width + x] = 0xFF666666 - add;
      else
        data[y * width + x] = 0xFFEEEEEE + add;
    }
  }

  munmap(data, size);

  add = add + 1;

  return buffer;
}

const static struct wl_callback_listener callback_listener;

static void render_buffer_to_surface(struct wl_client *client)
{
  struct wl_buffer *buffer = render_buffer(client->shm);
  wl_surface_attach(client->surface, buffer, 0, 0);

  struct wl_callback *wl_frame_callback = wl_surface_frame(client->surface);
  wl_callback_add_listener(wl_frame_callback, &callback_listener, client);
}

static void frame_done(void *data,
                       struct wl_callback *wl_callback,
                       uint32_t callback_data)
{
  printf("Frame done: %d\n", add);
  wl_callback_destroy(wl_callback);
  struct wl_client *client = data;
  render_buffer_to_surface(client);

  wl_surface_damage_buffer(client->surface, 0, 0, INT32_MAX, INT32_MAX);

  wl_surface_commit(client->surface);
}

const static struct wl_callback_listener callback_listener = {
    .done = frame_done,
};

static void xdg_surface_configure(void *data, struct xdg_surface *surface, uint32_t serial)
{
  struct wl_client *client = data;

  xdg_surface_ack_configure(surface, serial);
  render_buffer_to_surface(client);

  wl_surface_commit(client->surface);
}

const static struct xdg_surface_listener surface_listener = {
    .configure = xdg_surface_configure,
};

void add_and_render_surface(struct wl_client *client)
{
  client->surface = wl_compositor_create_surface(client->compositor);
  struct xdg_surface *xdg_surface = xdg_wm_base_get_xdg_surface(client->wm_base, client->surface);

  struct xdg_toplevel *xdg_toplevel = xdg_surface_get_toplevel(xdg_surface);
  xdg_toplevel_set_title(xdg_toplevel, "My f example");

  xdg_surface_add_listener(xdg_surface, &surface_listener, client);
  wl_surface_commit(client->surface);
}
