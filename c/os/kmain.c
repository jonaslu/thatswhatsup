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
  // This will set up segments and enter 32-bit protected mode
  init_gdt();

  fb_write_text("Printin digits:\n");
  fb_write_dec(1);
  fb_write_text("\n");
  fb_write_dec(0xFFFFFFFF);

  return 1;
}
