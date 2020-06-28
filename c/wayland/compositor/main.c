#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <wayland-server.h>
#include <time.h>

#include <wlr/backend.h>
#include <wlr/types/wlr_output.h>
#include <wlr/render/wlr_renderer.h>
#include <wlr/types/wlr_compositor.h>
#include <wlr/types/wlr_xdg_shell.h>

struct mcw_output
{
  struct wlr_output *wlr_output;
  struct mcw_server *server;

  struct wl_listener destroy;
  struct wl_listener frame;

  struct timespec last_frame;
  struct wl_list link;
};

struct mcw_server
{
  struct wl_display *wl_display;
  struct wl_event_loop *event_loop;

  struct wlr_backend *backend;
  struct wlr_compositor *compositor;

  struct wl_listener new_output;

  struct wl_list outputs;
};

static void frame_output_notify(struct wl_listener *listener, void *data)
{
  struct mcw_output *output = wl_container_of(listener, output, frame);
  struct wlr_output *wlr_output = data;

  struct wlr_renderer *renderer = wlr_backend_get_renderer(wlr_output->backend);

  wlr_output_attach_render(wlr_output, NULL);

  wlr_renderer_begin(renderer, wlr_output->width, wlr_output->height);
  float color[4] = {1.0, 0.0, 0.0, 1.0};
  wlr_renderer_clear(renderer, color);
  wlr_renderer_end(renderer);

  wlr_output_commit(wlr_output);
}

static void destroy_output_notify(struct wl_listener *listener, void *data)
{
  struct mcw_output *output = wl_container_of(listener, output, destroy);
  wl_list_remove(&output->link);
  wl_list_remove(&output->destroy.link);
  wl_list_remove(&output->frame.link);
  free(output);
}

static void new_output_notify(struct wl_listener *listener, void *data)
{
  struct mcw_server *server = wl_container_of(listener, server, new_output);
  struct wlr_output *wlr_output = data;

  if (!wl_list_empty(&wlr_output->modes))
  {
    struct wlr_output_mode *mode = wl_container_of(wlr_output->modes.prev, mode, link);
    wlr_output_set_mode(wlr_output, mode);
  }

  struct mcw_output *output = calloc(1, sizeof(struct mcw_output));
  clock_gettime(CLOCK_MONOTONIC, &output->last_frame);

  output->server = server;
  output->wlr_output = wlr_output;

  output->destroy.notify = destroy_output_notify;
  wl_signal_add(&wlr_output->events.destroy, &output->destroy);

  output->frame.notify = frame_output_notify;
  wl_signal_add(&wlr_output->events.frame, &output->frame);

  wlr_output_create_global(wlr_output);

  wl_list_insert(&server->outputs, &output->link);
}

int main()
{
  struct mcw_server server;

  server.wl_display = wl_display_create();
  assert(server.wl_display);

  server.event_loop = wl_display_get_event_loop(server.wl_display);
  assert(server.event_loop);

  server.backend = wlr_backend_autocreate(server.wl_display, NULL);
  assert(server.backend);

  wl_list_init(&server.outputs);

  server.new_output.notify = new_output_notify;
  wl_signal_add(&server.backend->events.new_output, &server.new_output);

  const char *socket = wl_display_add_socket_auto(server.wl_display);

  assert(socket);
  if (!socket) {
    printf("Could not create socket\n");
    return 1;
  }

  printf("Wayland running on socket: %s\n", socket);
  setenv("WAYLAND_DISPLAY", socket, true);

  wl_display_init_shm(server.wl_display);
  server.compositor = wlr_compositor_create(server.wl_display, wlr_backend_get_renderer(server.backend));
  wlr_xdg_shell_create(server.wl_display);

  if (!wlr_backend_start(server.backend))
  {
    printf("Failed to start server");
    return 1;
  }

  wl_display_run(server.wl_display);
  wl_display_destroy(server.wl_display);

  return 0;
}
