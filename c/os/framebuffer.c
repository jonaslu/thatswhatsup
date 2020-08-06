#include "io.h"
#include "framebuffer.h"

#define VIDEO_MEM 0x000B8000

#define FB_COMMAND_PORT 0x3D4
#define FB_DATA_PORT 0x3D5
#define FB_SHOW_CURSOR 0xA
#define FB_HIGH_BYTE_COMMAND 14
#define FB_LOW_BYTE_COMMAND 15

#define FB_COL_MAX 80
#define FB_ROW_MAX 25

unsigned int fb_row = 0;
unsigned int fb_column = 0;

unsigned char fb_fg = FB_WHITE;
unsigned char fb_bg = FB_BLACK;

unsigned char *fb = (unsigned char *)VIDEO_MEM;

static void fb_write_cell(unsigned int location, char c, unsigned char fg, unsigned char bg)
{
  unsigned int fb_pos = location * 2;

  fb[fb_pos] = c;
  fb[fb_pos + 1] = ((bg & 0x0F) << 4) | (fg & 0x0F);
}

static void memcpy(unsigned char *dest, unsigned char *src, unsigned int count)
{
  for (unsigned int i = 0; i < count; i++)
  {
    dest[i] = src[i];
  }
}

static void scroll_one_line()
{
  unsigned int startpos = FB_COL_MAX;

  unsigned int endpos = FB_ROW_MAX * FB_COL_MAX - 1;
  unsigned int count = endpos - startpos;

  memcpy(fb, fb + (startpos * 2), count * 2);

  for (unsigned int i = (FB_ROW_MAX - 1) * FB_COL_MAX - 1; i < endpos; i++)
  {
    fb_write_cell(i, ' ', fb_fg, fb_bg);
  }
}

void fb_write_text(const char *buf)
{
  for (unsigned int i = 0;; i++)
  {
    char next_char = buf[i];

    switch (next_char)
    {
    case '\0':
      return;
      break;

    case '\n':
      fb_column = 0;
      fb_row++;
      break;

    default:
    {
      unsigned int location = fb_row * FB_COL_MAX + fb_column;
      fb_write_cell(location, buf[i], fb_fg, fb_bg);
      fb_column++;
      break;
    }
    }

    if (fb_column >= FB_COL_MAX)
    {
      fb_column = 0;
      fb_row++;
    }

    if (fb_row >= FB_ROW_MAX)
    {
      scroll_one_line();
      fb_column = 0;
      fb_row = FB_ROW_MAX - 1;
    }

    unsigned int location = fb_row * FB_COL_MAX + fb_column;
    fb_move_cursor(location);
  }
}

void fb_set_fg_color(unsigned char fg)
{
  fb_fg = fg;
}

void fb_set_bg_color(unsigned char bg)
{
  fb_bg = bg;
}

void fb_show_cursor()
{
  outb(FB_COMMAND_PORT, FB_SHOW_CURSOR);
  outb(FB_DATA_PORT, 0x0);
}

void fb_move_cursor(unsigned int pos)
{
  outb(FB_COMMAND_PORT, FB_HIGH_BYTE_COMMAND);
  outb(FB_DATA_PORT, (pos >> 8) & 0x00FF);
  outb(FB_COMMAND_PORT, FB_LOW_BYTE_COMMAND);
  outb(FB_DATA_PORT, pos & 0x00FF);
}
