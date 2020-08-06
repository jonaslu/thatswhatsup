#include "framebuffer.h"

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

  char *buf = "12345\n";

  for (unsigned int i = 0; i < 25; i++)
  {
    fb_write_text(buf);
    (*buf)++;
  }

  fb_write_text("1235");

  return 1;
}
