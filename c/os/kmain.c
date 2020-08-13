#include "framebuffer.h"
#include "serial.h"
#include "gdt.h"

static void init_fb()
{
  fb_set_fg_color(FB_WHITE);
  fb_set_bg_color(FB_BLACK);

  fb_show_cursor();
  fb_move_cursor(0);
}

int kmain()
{
  init_fb();
  serial_init();
  init_gdt();

  fb_write_text("1235");

  return 1;
}
