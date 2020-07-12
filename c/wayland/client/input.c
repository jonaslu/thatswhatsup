#include <unistd.h>
#include <string.h>
#include <assert.h>
#include <sys/mman.h>
#include <xkbcommon/xkbcommon.h>
#include <wayland-client.h>

#include "client.h"

struct client_state
{
  struct wl_client *client;
  struct wl_keyboard *keyboard;
  struct xkb_context *xkb_context;
  struct xkb_keymap *xkb_keymap;
  struct xkb_state *xkb_state;
  char text[128];
  int text_index;
} state;

static void wl_keyboard_listener_keymap(void *data, struct wl_keyboard *wl_keyboard, uint32_t format, int32_t fd, uint32_t size)
{
  struct client_state *client_state = data;
  assert(format == WL_KEYBOARD_KEYMAP_FORMAT_XKB_V1);

  char *map_shm = mmap(NULL, size, PROT_READ, MAP_SHARED, fd, 0);
  assert(map_shm != MAP_FAILED);

  // OFC, it's here
  struct xkb_keymap *xkb_keymap = xkb_keymap_new_from_string(client_state->xkb_context, map_shm, XKB_KEYMAP_FORMAT_TEXT_V1, XKB_KEYMAP_COMPILE_NO_FLAGS);
  munmap(map_shm, size);
  close(fd);

  struct xkb_state *xkb_state = xkb_state_new(xkb_keymap);
  xkb_state_unref(client_state->xkb_state);
  xkb_keymap_unref(client_state->xkb_keymap);

  client_state->xkb_keymap = xkb_keymap;
  client_state->xkb_state = xkb_state;
}

static void wl_keyboard_listener_enter(void *data, struct wl_keyboard *wl_keyboard, uint32_t serial, struct wl_surface *surface, struct wl_array *keys)
{
}

static void wl_keyboard_listener_leave(void *data, struct wl_keyboard *wl_keyboard, uint32_t serial, struct wl_surface *surface)
{
}

static void wl_keyboard_listener_key(void *data, struct wl_keyboard *wl_keyboard, uint32_t serial, uint32_t time, uint32_t key, uint32_t state)
{
  if (state == WL_KEYBOARD_KEY_STATE_PRESSED)
  {
    struct client_state *client_state = data;
    char buf[128];
    uint8_t keycode = key + 8;

    xkb_keysym_t xkb_keysym = xkb_state_key_get_one_sym(client_state->xkb_state, keycode);
    if (xkb_keysym == XKB_KEY_BackSpace)
    {
      printf("Gettin here");
      client_state->text[strlen(client_state->text) - 1] = '\0';
    }
    else
    {
      int chars_read = xkb_state_key_get_utf8(client_state->xkb_state, keycode, buf, sizeof(buf));
      strcat(client_state->text, buf);
    }

    printf("Key pressed: %s\n", buf);
    render_buffer(client_state->client, client_state->text);
  }
}

static void wl_keyboard_listener_modifiers(void *data, struct wl_keyboard *wl_keyboard, uint32_t serial, uint32_t mods_depressed, uint32_t mods_latched, uint32_t mods_locked, uint32_t group)
{
}

static void wl_keyboard_listener_repeat_info(void *data, struct wl_keyboard *wl_keyboard, int32_t rate, int32_t delay)
{
}

const static struct wl_keyboard_listener wl_keyboard_listener = {
    .keymap = wl_keyboard_listener_keymap,
    .enter = wl_keyboard_listener_enter,
    .leave = wl_keyboard_listener_leave,
    .key = wl_keyboard_listener_key,
    .modifiers = wl_keyboard_listener_modifiers,
    .repeat_info = wl_keyboard_listener_repeat_info,
};

static void wl_seat_name(void *data, struct wl_seat *wl_seat, const char *name)
{
  printf("Wl seat name: %s\n", name);
}

static void wl_seat_capabilities(void *data,
                                 struct wl_seat *wl_seat,
                                 uint32_t capabilities)
{
  printf("Listener seat added\n");
  struct client_state *client_state = data;

  int have_keyboard = capabilities & WL_SEAT_CAPABILITY_KEYBOARD;

  if (have_keyboard && client_state->keyboard == NULL)
  {
    // printf("Client seat 2: %p\n", (void *)wl_seat);
    client_state->keyboard = wl_seat_get_keyboard(wl_seat);
    wl_keyboard_add_listener(client_state->keyboard, &wl_keyboard_listener, client_state);
  }
  else if (!have_keyboard && client_state->keyboard != NULL)
  {
    wl_keyboard_release(client_state->keyboard);
    client_state->keyboard = NULL;
  }
}

const static struct wl_seat_listener seat_listener = {
    .capabilities = wl_seat_capabilities,
    .name = wl_seat_name,
};

void init_kbd_input(struct wl_client *client)
{
  state.xkb_context = xkb_context_new(XKB_CONTEXT_NO_FLAGS);
  // state.xkb_keymap = xkb_keymap_new_from_string(state.xkb_context, "se", XKB_KEYMAP_FORMAT_TEXT_V1, XKB_KEYMAP_COMPILE_NO_FLAGS);
  // state.xkb_state = xkb_state_new(state.xkb_keymap);
  // printf("Client seat: %p\n", (void *)client->seat);
  state.client = client;
  wl_seat_add_listener(client->seat, &seat_listener, &state);
}
