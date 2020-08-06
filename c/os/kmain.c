#include "framebuffer.h"

int kmain()
{
  fb_set_fg_color(FB_WHITE);
  fb_set_bg_color(FB_BLACK);

  fb_reset_cursor(0);

  char *buf = "1\n";

  for (unsigned int i = 0; i < 27; i++)
  {
    fb_write_text(buf);
    (*buf)++;
  }

  return 1;
}
