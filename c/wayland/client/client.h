#ifndef __WL_CLIENT
#define __WL_CLIENT

struct wl_client
{
  struct wl_display *display;
  struct wl_compositor *compositor;
  struct wl_registry *registry;
  struct xdg_wm_base *wm_base;
  struct wl_shm *shm;
  struct wl_surface *surface;
  struct wl_seat *seat;
};

void init_registry_listener(struct wl_client *client);
void init_kbd_input(struct wl_client *client);
void init_surface(struct wl_client *client);
void init_kbd_input(struct wl_client *client);
void render_chars(struct wl_client *client, const char *text);

int init_pty();
void read_from_pty(struct wl_client *wl_client);
void write_to_pty(const char *text, int size);

#endif
