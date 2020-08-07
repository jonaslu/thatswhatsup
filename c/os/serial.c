#include "io.h"

#define SERIAL_COM1_BASE 0x3F8
#define SERIAL_DATA_PORT(base) (base)
#define SERIAL_FIFO_COMMAND_PORT(base) (base + 2)
#define SERIAL_LINE_COMMAND_PORT(base) (base + 3)
#define SERIAL_MODEM_COMMAND_PORT(base) (base + 4)
#define SERIAL_LINE_STATUS_PORT(base) (base + 5)

/* Mode setting: expect the highest 8 bits first, then the low 8 bits */
#define SERIAL_LINE_ENABLE_DLAB 0x80
/* Length of 8 bits, no parity bits, one stop bit and break control disabled */
#define SERIAL_LINE_CONFIGURATION 0x03

#define SERIAL_FIFO_CONFIGURATION 0b11000111

#define SERIAL_MODEM_CONFIGURATION 0x03

// Send the baud rate = a divisor of 115200 bits/s
static void serial_configure_baud_rate(unsigned short com, unsigned short divisor)
{
  outb(SERIAL_LINE_COMMAND_PORT(com), SERIAL_LINE_ENABLE_DLAB);
  outb(SERIAL_DATA_PORT(com), (divisor >> 8) && 0x00FF);
  outb(SERIAL_DATA_PORT(com), divisor && 0x00FF);
}

// Set bits regarding the data we're about to send
static void serial_configure_line(unsigned short com)
{
  outb(SERIAL_LINE_COMMAND_PORT(com), SERIAL_LINE_CONFIGURATION);
}

static void serial_configure_fifo(unsigned short com)
{
  outb(SERIAL_FIFO_COMMAND_PORT(com), SERIAL_FIFO_CONFIGURATION);
}

static void serial_configure_modem(unsigned short com)
{
  outb(SERIAL_MODEM_COMMAND_PORT(com), SERIAL_MODEM_CONFIGURATION);
}

static int serial_is_transmit_fifo_empty(unsigned short com)
{
  return inb(SERIAL_LINE_STATUS_PORT(com)) & 0x20;
}

void serial_init() {
  serial_configure_baud_rate(SERIAL_COM1_BASE, 1);

  serial_configure_line(SERIAL_COM1_BASE);
  serial_configure_fifo(SERIAL_COM1_BASE);
  serial_configure_modem(SERIAL_COM1_BASE);
}

void serial_write(char *buffer)
{
  for (int i = 0;; i++)
  {
    char next_char = buffer[i];

    if (next_char == '\0')
    {
      return;
    }

    while (!serial_is_transmit_fifo_empty(SERIAL_COM1_BASE)) {};

    outb(SERIAL_DATA_PORT(SERIAL_COM1_BASE), next_char);
  }
}
