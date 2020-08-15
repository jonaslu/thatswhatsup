#include "framebuffer.h"
#include "serial.h"
#include "gdt.h"
#include "idt.h"

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
  init_idt();

  asm volatile ("int $0x3");
  asm volatile ("int $0x4");

  for(;;);

  return 1;
}
