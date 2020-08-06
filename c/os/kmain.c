#include "framebuffer.h"

int kmain()
{
  char *buf = "1\n";

  for (unsigned int i = 0; i < 27; i++)
  {
    write(buf);
    (*buf)++;
  }

  return 1;
}
