#define VIDEO_MEM 0x000B8000

#define FB_BLACK 0
#define FB_BLUE 1
#define FB_GREEN 2
#define FB_CYAN 3
#define FB_RED 4
#define FB_MAGENTA 5
#define FB_BROWN 6
#define FB_LIGHT_GREY 7
#define FB_DARK_GREY 8
#define FB_LIGHT_BLUE 9
#define FB_LIGHT_GREEN 10
#define FB_LIGHT_CYAN 11
#define FB_LIGHT_RED 12
#define FB_LIGHT_MAGENTA 13
#define FB_LIGHT_BROWN 14
#define FB_WHITE 15

void fb_write_cell(unsigned int location, char c, unsigned char fg, unsigned char bg)
{
  char *fb = (char *)VIDEO_MEM;

  fb[location] = c;
  fb[location + 1] = ((bg & 0x0F) << 4) | (fg & 0x0F);
}

int kmain()
{
  fb_write_cell(0, 'B', FB_LIGHT_GREEN, FB_BLACK);
  fb_write_cell(2, 'C', FB_LIGHT_GREEN, FB_BLACK);

  return 1;
}
