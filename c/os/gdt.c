#include "io.h"

struct gdt_entry
{
  unsigned short limit_low;
  unsigned short base_low;
  unsigned char base_middle;
  unsigned char access;
  unsigned char granularity;
  unsigned char base_high;
} __attribute__((packed));

struct gdt_ptr_t
{
  unsigned short size;
  unsigned int address;
} __attribute__((packed));

// Define this globally so it won't dissapear on the stack
struct gdt_entry gdt[3];
struct gdt_ptr_t gdt_ptr;

static void init_gdt_entry(unsigned char index, unsigned int base, unsigned int limit, unsigned char access, unsigned char granularity)
{
  gdt[index].base_low = (base & 0xFFFF);
  gdt[index].base_middle = (base >> 16) & 0xFF;
  gdt[index].base_high = (base >> 24) & 0xFF;

  gdt[index].limit_low = (limit & 0xFFFF);
  gdt[index].granularity = (limit >> 16) & 0x0F;

  gdt[index].granularity  |= granularity & 0xF0;
  gdt[index].access = access;
}

void init_gdt()
{
  init_gdt_entry(0, 0, 0, 0, 0);
  init_gdt_entry(1, 0, 0xFFFFFFFF, 0x9A, 0xCF);
  init_gdt_entry(2, 0, 0xFFFFFFFF, 0x92, 0xCF);

  gdt_ptr.size = (sizeof(struct gdt_entry) * 3) - 1;
  gdt_ptr.address = (unsigned int) &gdt;

  lgdt((unsigned int) &gdt_ptr);
}
